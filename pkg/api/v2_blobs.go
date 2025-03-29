package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/helpers"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
	"github.com/oklog/ulid/v2"
	"golang.org/x/sync/errgroup"
)

const (
	_ = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
)

type V2BlobsHandler struct {
	Config              *config.AppConfig
	S3Client            *s3.Client
	S3PresignClient     *s3.PresignClient
	BlobChunkRepository *repositories.BlobChunkRepository
}

// This handles the initiation of an upload session and monolithic upload
func (h *V2BlobsHandler) InitiateUploadSession(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("repositoryName").(string)
	org := r.Context().Value("organization").(*ent.Organization)
	registry := r.Context().Value("registry").(*ent.Registry)
	uploadId := ulid.Make().String()
	blobDigest := r.URL.Query().Get("digest")
	mount := r.URL.Query().Get("mount")

	// We check if this request is asking us to mount a blob from a different repository
	if err := h.mountBlob(r.Context(), org.Slug, mount); err == nil {
		w.Header().Set("Location", getBlobDownloadUrl(h.Config.GetBaseUrl(), org.Slug, registry.Slug, imageName, mount))
		w.Header().Set("Docker-Content-Digest", mount)
		w.WriteHeader(http.StatusCreated)
		return
	}

	// Create an S3 multipart upload
	output, err := h.S3Client.CreateMultipartUpload(r.Context(), &s3.CreateMultipartUploadInput{
		Bucket: &h.Config.S3.BlobsBucketName,
		Key:    aws.String(h.getKeyForBlobInFlight(org.Slug, uploadId)),
	})
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	sessionId := *output.UploadId

	contentType := r.Header.Get("Content-Type")

	var fullBytesRead int = -1

	// Check if this is a monolith upload, then handle the streaming upload to S3
	if blobDigest != "" && h.isMonolithicUpload(r.ContentLength, contentType) {
		keyName := h.getKeyForBlobInFlight(org.Slug, uploadId)

		fullBytesRead, err = h.handleStreamingUpload(r.Context(), r.Body, keyName, sessionId)
		if err != nil {
			responses.OCIBlobUploadInvalid(w)
			return
		}

		err = h.completeMultiPartUpload(r.Context(), org.Slug, keyName, sessionId, imageName, blobDigest)
		if err != nil {
			responses.OCIBlobUploadInvalid(w)
			return
		}
	}

	status := http.StatusAccepted

	if fullBytesRead > 0 {
		// If we recieved data then return the location of the blob and range of the uploaded data
		w.Header().Set("Range", fmt.Sprintf("0-%d", fullBytesRead))
		w.Header().Set("Location", getBlobDownloadUrl(h.Config.GetBaseUrl(), org.Slug, registry.Slug, imageName, blobDigest))
		status = http.StatusCreated
	} else {
		// Return the location for the upload and the minimum size for a chunk
		w.Header().Set("Location", h.getUploadUrl(org.Slug, registry.Slug, imageName, uploadId, sessionId))
		w.Header().Set("OCI-Chunk-Min-Length", strconv.Itoa(int(h.Config.ChunkMinLength)))
	}

	w.WriteHeader(status)
}

// If the client is requesting us to mount a blob from a different repository then we check in S3 if it exists
func (h *V2BlobsHandler) mountBlob(ctx context.Context, orgSlug string, digest string) error {
	if digest == "" {
		return errors.New("digest is empty")
	}

	_, found, err := h.s3HeadBlob(ctx, orgSlug, digest)
	if !found {
		return errors.New("blob not found")
	}
	if err != nil {
		return err
	}

	return nil
}

// This registry doesn't support blob deletion, so we return method not allowed
// this is still OCI compliant
func (h *V2BlobsHandler) DeleteBlob(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// This handler handles the HEAD method to check if a blob exists in the repository
func (h *V2BlobsHandler) HeadBlob(w http.ResponseWriter, r *http.Request) {
	blobDigest := r.Context().Value("digest").(string)
	org := r.Context().Value("organization").(*ent.Organization)

	output, found, err := h.s3HeadBlob(r.Context(), org.Slug, blobDigest)
	if !found {
		responses.OCIBlobUnknown(w, blobDigest)
		return
	}
	if err != nil {
		responses.OCIInternalServerError(w)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(int(*output.ContentLength)))
	w.Header().Set("Docker-Content-Digest", blobDigest)
	w.WriteHeader(http.StatusOK)
}

// This handles the download of a blob, it generates an S3 presign linked
// and redirects the user directly to S3 for download
func (h *V2BlobsHandler) DownloadBlob(w http.ResponseWriter, r *http.Request) {
	blobDigest := r.Context().Value("digest").(string)
	org := r.Context().Value("organization").(*ent.Organization)
	keyName := helpers.GetS3KeyForBlob(org.Slug, blobDigest)

	withFileName := r.URL.Query().Get("filename")

	getObjectInput := &s3.GetObjectInput{
		Bucket: &h.Config.S3.BlobsBucketName,
		Key:    &keyName,
	}

	if withFileName != "" {
		getObjectInput.ResponseContentDisposition = aws.String(fmt.Sprintf("attachment; filename=\"%s\"", withFileName))
	}

	req, err := h.S3PresignClient.PresignGetObject(r.Context(), getObjectInput, s3.WithPresignExpires(*aws.Duration(time.Minute * 10)))
	if err != nil {
		responses.OCIInternalServerError(w)
		log.Println(err)
		return
	}

	w.Header().Set("Location", req.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Handles the upload of a blob
func (h *V2BlobsHandler) UploadBlob(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("repositoryName").(string)
	uploadId := r.Context().Value("uploadId").(string)
	org := r.Context().Value("organization").(*ent.Organization)
	registry := r.Context().Value("registry").(*ent.Registry)
	sessionId := r.URL.Query().Get("session")
	defer r.Body.Close()

	if sessionId == "" {
		responses.OCIBlobUploadUnknown(w)
		return
	}

	keyName := h.getKeyForBlobInFlight(org.Slug, uploadId)

	// Parse and validate the content-range header to determine if this is a
	// chunked upload
	contentRange := r.Header.Get("Content-Range")
	rangeFrom, rangeTo, err := h.parseContentRange(contentRange)
	if err != nil {
		responses.OCIBlobUploadInvalid(w)
		return
	}
	isChunkUpload := rangeTo > 0

	if isChunkUpload {
		// Chunks have to be uploaded in order, check that the upload is valid
		isOutOfOrder, err := h.BlobChunkRepository.IsOutOfOrder(r.Context(), sessionId, uploadId, rangeFrom, rangeTo)
		if err != nil {
			responses.OCIBlobUploadInvalid(w)
			return
		}

		if isOutOfOrder {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		// Get the next chunk that is going to be uploaded
		chunk, err := h.BlobChunkRepository.GetNext(r.Context(), sessionId, uploadId, rangeFrom, rangeTo)
		if err != nil {
			log.Println(err)
			responses.OCIBlobUploadInvalid(w)
			return
		}

		// Build the upload request to S3
		partInput := &s3.UploadPartInput{
			ContentLength: &r.ContentLength,
			Body:          r.Body,
			Bucket:        &h.Config.S3.BlobsBucketName,
			Key:           &keyName,
			PartNumber:    aws.Int32(int32(chunk.PartNumber)),
			UploadId:      &sessionId,
		}

		// Upload each chunk to S3
		_, err = h.S3Client.UploadPart(r.Context(), partInput)
		if err != nil {
			log.Println(err)
			responses.OCIBlobUploadInvalid(w)
			return
		}
	} else {
		// If this is not a chunked uploaded then handle it as a streaming upload and chunk on the fly
		fullBytesRead, err := h.handleStreamingUpload(r.Context(), r.Body, keyName, sessionId)
		if err != nil {
			responses.OCIBlobUploadInvalid(w)
			return
		}
		rangeTo = uint64(fullBytesRead)
	}

	// Return back the uploaded range of the blob
	w.Header().Set("Range", fmt.Sprintf("%d-%d", rangeFrom, rangeTo))
	w.Header().Set("Location", h.getUploadUrl(org.Slug, registry.Slug, imageName, uploadId, sessionId))
	w.WriteHeader(http.StatusAccepted)
}

// This handler finalizes an upload session by completing the multi part upload in S3, it also
// supports receiving a final upload chunk as part of the request.
func (h *V2BlobsHandler) FinalizeBlobUploadSession(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("repositoryName").(string)
	uploadId := r.Context().Value("uploadId").(string)
	org := r.Context().Value("organization").(*ent.Organization)
	registry := r.Context().Value("registry").(*ent.Registry)
	sessionId := r.URL.Query().Get("session")

	if sessionId == "" {
		responses.OCIBlobUploadUnknown(w)
		return
	}

	// If the blob's digest is not in the request then we can't finalize it
	blobDigest := r.URL.Query().Get("digest")
	if blobDigest == "" {
		responses.OCIBlobUploadInvalid(w)
		return
	}

	keyName := h.getKeyForBlobInFlight(org.Slug, uploadId)

	contentType := r.Header.Get("Content-Type")

	var err error

	// Check if we are receving a final part
	if h.isMonolithicUpload(r.ContentLength, contentType) {
		_, err = h.handleStreamingUpload(r.Context(), r.Body, keyName, sessionId)
		if err != nil {
			responses.OCIBlobUploadInvalid(w)
			return
		}
	}

	// Complete the multipart upload in S3
	err = h.completeMultiPartUpload(r.Context(), org.Slug, keyName, sessionId, imageName, blobDigest)
	if err != nil {
		log.Println(err)
		responses.OCIBlobUploadInvalid(w)
		return
	}

	// Clean up the temporal data in the DB
	err = h.BlobChunkRepository.DeleteAllForUploadID(r.Context(), uploadId)
	if err != nil {
		responses.OCIInternalServerError(w)
		return
	}

	// Return back the blob's digest and its location for download
	w.Header().Set("Docker-Content-Digest", blobDigest)
	w.Header().Set("Location", getBlobDownloadUrl(h.Config.GetBaseUrl(), org.Slug, registry.Slug, imageName, blobDigest))
	w.WriteHeader(http.StatusCreated)
}

// This handles retreving an existing upload session for a blob
func (h *V2BlobsHandler) GetBlobUploadSession(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("repositoryName").(string)
	uploadId := r.Context().Value("uploadId").(string)
	org := r.Context().Value("organization").(*ent.Organization)
	registry := r.Context().Value("registry").(*ent.Registry)

	// Check the database for the upload session
	chunks, err := h.BlobChunkRepository.GetByUploadID(r.Context(), uploadId)
	if err != nil {
		responses.OCIInternalServerError(w)
		return
	}

	if len(chunks) == 0 {
		responses.OCIBlobUploadUnknown(w)
		return
	}

	totalUploadedRange := 0

	// Calculate the total size of all the uploaded chunks
	// so the client knows how to resume the upload
	for _, chunk := range chunks {
		totalUploadedRange += int(chunk.RangeTo)
	}

	w.Header().Set("Range", fmt.Sprintf("0-%d", totalUploadedRange))
	w.Header().Set("Location", h.getUploadUrl(org.Slug, registry.Slug, imageName, uploadId, chunks[0].SessionID))
	w.WriteHeader(http.StatusNoContent)
}

// Function used to async upload a part to S3
func (h *V2BlobsHandler) asyncPartUpload(ctx context.Context, objectKey string, buffer []byte, bytesRead int, partNumber int, uploadId string) (*types.CompletedPart, error) {
	partInput := &s3.UploadPartInput{
		Body:       manager.ReadSeekCloser(bytes.NewReader(buffer[:bytesRead])),
		Bucket:     &h.Config.S3.BlobsBucketName,
		Key:        &objectKey,
		PartNumber: aws.Int32(int32(partNumber)),
		UploadId:   &uploadId,
	}

	log.Println("uploading partNumber:", *partInput.PartNumber)

	// Upload each chunk to S3
	uploadResult, err := h.S3Client.UploadPart(ctx, partInput)
	if err != nil {
		log.Println("error uploading partNumber: ", partNumber, err.Error())
		return nil, fmt.Errorf("failed to upload part to S3: %v", err.Error())
	}

	log.Println("uploading partNumber: ", partNumber, " finished")

	return &types.CompletedPart{
		ETag:       uploadResult.ETag,
		PartNumber: aws.Int32(int32(partNumber)),
	}, nil
}

// This function handlers a streaming upload from a client. As the data is being streamed from the client
// it chunks up the request and stream uploads it to S3 in parts, this is done to handle an edge case
// where certain clients upload big blobs without utlising the OCI blob chunk flow specified in the specs
func (h *V2BlobsHandler) handleStreamingUpload(ctx context.Context, body io.ReadCloser, keyName string, sessionId string) (int, error) {
	partNumber := 0
	fullBytesRead := 0

	eg, grpCtx := errgroup.WithContext(ctx)
	// Set the maximum amount of go routines to run at once
	eg.SetLimit(h.Config.BlobUploadMaxGoRoutines)

	// Read file in chunks
	for {
		// Read chunk to buffer
		buffer := make([]byte, h.Config.ChunkBufferLength)
		bytesRead, readErr := io.ReadFull(body, buffer)

		if bytesRead > 0 {
			partNumber++
			// Async upload to S3
			eg.Go(func() error {
				_, err := h.asyncPartUpload(grpCtx, keyName, buffer, bytesRead, partNumber, sessionId)
				if err != nil {
					return err
				}
				return nil
			})
			fullBytesRead += bytesRead
		}

		// Break out of the loop if end of file is reached
		if readErr == io.EOF || readErr == io.ErrUnexpectedEOF {
			break
		} else if readErr != nil {
			return 0, errors.New("failed to read the file")
		}
	}

	if err := eg.Wait(); err != nil {
		return 0, errors.New("failed to upload part to S3")
	}

	return fullBytesRead, nil
}

// Complete an S3 multipart upload to finish an upload
func (h *V2BlobsHandler) completeMultiPartUpload(ctx context.Context, orgSlug string, keyName string, sessionId string, imageName string, blobDigest string) error {
	// Get all the uploaded parts from S3
	partsOutput, err := h.S3Client.ListParts(ctx, &s3.ListPartsInput{
		Bucket:   &h.Config.S3.BlobsBucketName,
		Key:      &keyName,
		UploadId: &sessionId,
	})
	if err != nil {
		return err
	}

	var completedParts []types.CompletedPart

	for _, p := range partsOutput.Parts {
		completedParts = append(completedParts, types.CompletedPart{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
		})
	}

	// Complete the multi part upload
	_, err = h.S3Client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   &h.Config.S3.BlobsBucketName,
		Key:      &keyName,
		UploadId: &sessionId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		return err
	}

	// Check the blob exists to retrieve some additional metadata
	headResp, err := h.S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &h.Config.S3.BlobsBucketName,
		Key:    &keyName,
	})
	if err != nil {
		return err
	}

	objectSize := *headResp.ContentLength

	if objectSize > (5 * GiB) {
		return errors.New("object size is above 5GB")
	}

	destKey := helpers.GetS3KeyForBlob(orgSlug, blobDigest)
	copySource := fmt.Sprintf("%s/%s", h.Config.S3.BlobsBucketName, keyName)

	// Copy the in-flight blob to its final path in S3
	_, err = h.S3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     &h.Config.S3.BlobsBucketName,
		CopySource: &copySource,
		Key:        &destKey,
	})

	return err
}

// Check in S3 if a blob exists
func (h *V2BlobsHandler) s3HeadBlob(ctx context.Context, orgSlug string, blobDigest string) (*s3.HeadObjectOutput, bool, error) {
	output, err := h.S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &h.Config.S3.BlobsBucketName,
		Key:    aws.String(helpers.GetS3KeyForBlob(orgSlug, blobDigest)),
	})

	var nfe *types.NotFound
	if err != nil && errors.As(err, &nfe) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	return output, true, nil
}

// Generate the S3 key for a blob that is "in-flight" meaning that it hasn't been fully uploaded
func (h *V2BlobsHandler) getKeyForBlobInFlight(orgSlug string, uploadId string) string {
	return fmt.Sprintf("%s/in-flight/%s.blob", orgSlug, uploadId)
}

// Generate the URL for uploading a blob
func (h *V2BlobsHandler) getUploadUrl(orgSlug string, registrySlug string, imageName string, uploadId string, s3UploadId string) string {
	return fmt.Sprintf("%s/v2/%s/%s/%s/blobs/uploads/%s?session=%s", h.Config.GetBaseUrl(), orgSlug, registrySlug, imageName, uploadId, s3UploadId)
}

// Generate the URL for downloading a blob
func getBlobDownloadUrl(baseURL string, orgSlug string, registrySlug string, imageName string, digest string) string {
	return fmt.Sprintf("%s/v2/%s/%s/%s/blobs/%s", baseURL, orgSlug, registrySlug, imageName, digest)
}

// Check if this is a monolithic upload
func (h *V2BlobsHandler) isMonolithicUpload(contentLength int64, contentType string) bool {
	return contentLength > 0 && contentType == "application/octet-stream"
}

// Parse the content range of an upload request and check that is valid
func (h *V2BlobsHandler) parseContentRange(contentRange string) (uint64, uint64, error) {
	if contentRange == "" {
		return 0, 0, nil
	}

	ranges := strings.Split(contentRange, "-")
	if len(ranges) < 2 {
		return 0, 0, errors.New("invalid content range")
	}

	rangeFrom, err := strconv.ParseUint(ranges[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	rangeTo, err := strconv.ParseUint(ranges[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return rangeFrom, rangeTo, nil
}
