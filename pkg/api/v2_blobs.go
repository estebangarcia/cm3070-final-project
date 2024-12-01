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

func (h *V2BlobsHandler) InitiateUploadSession(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	uploadId := ulid.Make().String()
	blobDigest := r.URL.Query().Get("digest")
	mount := r.URL.Query().Get("mount")

	if err := h.mountBlob(r.Context(), mount); err == nil {
		w.Header().Set("Location", h.getBlobDownloadUrl(imageName, mount))
		w.Header().Set("Docker-Content-Digest", mount)
		w.WriteHeader(http.StatusCreated)
		return
	}

	output, err := h.S3Client.CreateMultipartUpload(r.Context(), &s3.CreateMultipartUploadInput{
		Bucket: &h.Config.S3.BlobsBucketName,
		Key:    aws.String(h.getKeyForBlobInFlight(uploadId)),
	})
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	sessionId := *output.UploadId

	contentType := r.Header.Get("Content-Type")

	var fullBytesRead int = -1

	if blobDigest != "" && h.isMonolithicUpload(r.ContentLength, contentType) {
		keyName := h.getKeyForBlobInFlight(uploadId)

		fullBytesRead, err = h.handleStreamingUpload(r.Context(), r.Body, keyName, sessionId)
		if err != nil {
			responses.OCIBlobUploadInvalid(w)
			return
		}

		err = h.completeMultiPartUpload(r.Context(), keyName, sessionId, imageName, blobDigest)
		if err != nil {
			responses.OCIBlobUploadInvalid(w)
			return
		}
	}

	status := http.StatusAccepted

	if fullBytesRead > 0 {
		w.Header().Set("Range", fmt.Sprintf("0-%d", fullBytesRead))
		w.Header().Set("Location", h.getBlobDownloadUrl(imageName, blobDigest))
		status = http.StatusCreated
	} else {
		w.Header().Set("Location", h.getUploadUrl(imageName, uploadId, sessionId))
		w.Header().Set("OCI-Chunk-Min-Length", strconv.Itoa(int(h.Config.ChunkMinLength)))
	}

	w.WriteHeader(status)
}

func (h *V2BlobsHandler) mountBlob(ctx context.Context, digest string) error {
	if digest == "" {
		return errors.New("digest is empty")
	}

	_, found, err := h.s3HeadBlob(ctx, digest)
	if !found {
		return errors.New("blob not found")
	}
	if err != nil {
		return err
	}

	return nil
}

func (h *V2BlobsHandler) HeadBlob(w http.ResponseWriter, r *http.Request) {
	blobDigest := r.Context().Value("digest").(string)

	output, found, err := h.s3HeadBlob(r.Context(), blobDigest)
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

func (h *V2BlobsHandler) DownloadBlob(w http.ResponseWriter, r *http.Request) {
	blobDigest := r.Context().Value("digest").(string)

	keyName := h.getKeyForBlob(blobDigest)

	req, err := h.S3PresignClient.PresignGetObject(r.Context(), &s3.GetObjectInput{
		Bucket: &h.Config.S3.BlobsBucketName,
		Key:    &keyName,
	}, s3.WithPresignExpires(*aws.Duration(time.Minute * 10)))
	if err != nil {
		responses.OCIInternalServerError(w)
		log.Println(err)
		return
	}

	w.Header().Set("Location", req.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *V2BlobsHandler) UploadBlob(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	uploadId := r.Context().Value("uploadId").(string)
	sessionId := r.URL.Query().Get("session")
	defer r.Body.Close()

	if sessionId == "" {
		responses.OCIBlobUploadUnknown(w)
		return
	}

	keyName := h.getKeyForBlobInFlight(uploadId)

	contentRange := r.Header.Get("Content-Range")
	rangeFrom, rangeTo, err := h.parseContentRange(contentRange)
	if err != nil {
		responses.OCIBlobUploadInvalid(w)
		return
	}
	isChunkUpload := rangeTo > 0

	if isChunkUpload {
		isOutOfOrder, err := h.BlobChunkRepository.IsOutOfOrder(r.Context(), sessionId, uploadId, rangeFrom, rangeTo)
		if err != nil {
			responses.OCIBlobUploadInvalid(w)
			return
		}

		if isOutOfOrder {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		chunk, err := h.BlobChunkRepository.GetNext(r.Context(), sessionId, uploadId, rangeFrom, rangeTo)
		if err != nil {
			log.Println(err)
			responses.OCIBlobUploadInvalid(w)
			return
		}

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
		fullBytesRead, err := h.handleStreamingUpload(r.Context(), r.Body, keyName, sessionId)
		if err != nil {
			responses.OCIBlobUploadInvalid(w)
			return
		}
		rangeTo = uint64(fullBytesRead)
	}

	w.Header().Set("Range", fmt.Sprintf("%d-%d", rangeFrom, rangeTo))
	w.Header().Set("Location", h.getUploadUrl(imageName, uploadId, sessionId))
	w.WriteHeader(http.StatusAccepted)
}

func (h *V2BlobsHandler) FinalizeBlobUploadSession(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	uploadId := r.Context().Value("uploadId").(string)
	sessionId := r.URL.Query().Get("session")

	if sessionId == "" {
		responses.OCIBlobUploadUnknown(w)
		return
	}

	blobDigest := r.URL.Query().Get("digest")
	if sessionId == "" {
		responses.OCIBlobUploadInvalid(w)
		return
	}

	keyName := h.getKeyForBlobInFlight(uploadId)

	contentType := r.Header.Get("Content-Type")

	var fullBytesRead int = -1
	var err error

	if h.isMonolithicUpload(r.ContentLength, contentType) {
		fullBytesRead, err = h.handleStreamingUpload(r.Context(), r.Body, keyName, sessionId)
		if err != nil {
			responses.OCIBlobUploadInvalid(w)
			return
		}
	}

	err = h.completeMultiPartUpload(r.Context(), keyName, sessionId, imageName, blobDigest)
	if err != nil {
		log.Println(err)
		responses.OCIBlobUploadInvalid(w)
		return
	}

	err = h.BlobChunkRepository.DeleteAllForUploadID(r.Context(), uploadId)
	if err != nil {
		responses.OCIInternalServerError(w)
		return
	}

	if fullBytesRead > 0 {
		w.Header().Set("Range", fmt.Sprintf("0-%d", fullBytesRead))
	}
	w.Header().Set("Location", h.getBlobDownloadUrl(imageName, blobDigest))
	w.WriteHeader(http.StatusCreated)
}

func (h *V2BlobsHandler) GetBlobUploadSession(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	uploadId := r.Context().Value("uploadId").(string)

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

	for _, chunk := range chunks {
		totalUploadedRange += int(chunk.RangeTo)
	}

	w.Header().Set("Range", fmt.Sprintf("0-%d", totalUploadedRange))
	w.Header().Set("Location", h.getUploadUrl(imageName, uploadId, chunks[0].SessionID))
	w.WriteHeader(http.StatusNoContent)
}

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

func (h *V2BlobsHandler) handleStreamingUpload(ctx context.Context, body io.ReadCloser, keyName string, sessionId string) (int, error) {
	partNumber := 0
	fullBytesRead := 0

	eg := errgroup.Group{}
	eg.SetLimit(h.Config.BlobUploadMaxGoRoutines)

	// Read file in chunks
	for {
		buffer := make([]byte, h.Config.ChunkBufferLength)
		bytesRead, readErr := io.ReadFull(body, buffer)

		if bytesRead > 0 {
			partNumber++
			eg.Go(func() error {
				_, err := h.asyncPartUpload(ctx, keyName, buffer, bytesRead, partNumber, sessionId)
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

func (h *V2BlobsHandler) completeMultiPartUpload(ctx context.Context, keyName string, sessionId string, imageName string, blobDigest string) error {
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

	destKey := h.getKeyForBlob(blobDigest)
	copySource := fmt.Sprintf("%s/%s", h.Config.S3.BlobsBucketName, keyName)

	_, err = h.S3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     &h.Config.S3.BlobsBucketName,
		CopySource: &copySource,
		Key:        &destKey,
	})

	return err
}

func (h *V2BlobsHandler) s3HeadBlob(ctx context.Context, blobDigest string) (*s3.HeadObjectOutput, bool, error) {
	output, err := h.S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &h.Config.S3.BlobsBucketName,
		Key:    aws.String(h.getKeyForBlob(blobDigest)),
	})

	var nfe *types.NotFound
	if err != nil && errors.As(err, &nfe) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	return output, true, nil
}

func (h *V2BlobsHandler) getKeyForBlobInFlight(uploadId string) string {
	return fmt.Sprintf("in-flight/%s.blob", uploadId)
}

func (h *V2BlobsHandler) getKeyForBlob(digest string) string {
	return fmt.Sprintf("blobs/%s/blob.data", helpers.GetDigestAsNestedFolder(digest))
}

func (h *V2BlobsHandler) getUploadUrl(imageName string, uploadId string, s3UploadId string) string {
	return fmt.Sprintf("%s/v2/%s/blobs/uploads/%s?session=%s", h.Config.BaseURL, imageName, uploadId, s3UploadId)
}

func (h *V2BlobsHandler) getBlobDownloadUrl(imageName string, digest string) string {
	return fmt.Sprintf("%s/v2/%s/blobs/%s", h.Config.BaseURL, imageName, digest)
}

func (h *V2BlobsHandler) isMonolithicUpload(contentLength int64, contentType string) bool {
	return contentLength > 0 && contentType == "application/octet-stream"
}

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
