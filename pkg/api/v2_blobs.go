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

const sha256Prefix = "sha256:"

type V2BlobsHandler struct {
	Config          *config.AppConfig
	S3Client        *s3.Client
	S3PresignClient *s3.PresignClient
}

func (h *V2BlobsHandler) InitiateUploadSession(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	uploadId := ulid.Make().String()

	output, err := h.S3Client.CreateMultipartUpload(r.Context(), &s3.CreateMultipartUploadInput{
		Bucket: &h.Config.S3.BlobsBucketName,
		Key:    aws.String(h.getKeyForBlobInFlight(imageName, uploadId)),
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", h.getUploadUrl(imageName, uploadId, *output.UploadId))
	w.Header().Set("OCI-Chunk-Min-Length", strconv.Itoa(int(h.Config.ChunkMinLength)))
	w.WriteHeader(http.StatusAccepted)
}

func (h *V2BlobsHandler) HeadBlob(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	blobDigest := r.Context().Value("digest").(string)

	output, err := h.S3Client.HeadObject(r.Context(), &s3.HeadObjectInput{
		Bucket: &h.Config.S3.BlobsBucketName,
		Key:    aws.String(h.getKeyForBlob(imageName, blobDigest)),
	})

	if err != nil {
		var nfe *types.NotFound
		if errors.As(err, &nfe) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(int(*output.ContentLength)))
	w.Header().Set("Docker-Content-Digest", blobDigest)
	w.WriteHeader(http.StatusOK)
}

func (h *V2BlobsHandler) DownloadBlob(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	blobDigest := r.Context().Value("digest").(string)

	keyName := h.getKeyForBlob(imageName, blobDigest)

	req, err := h.S3PresignClient.PresignGetObject(r.Context(), &s3.GetObjectInput{
		Bucket: &h.Config.S3.BlobsBucketName,
		Key:    &keyName,
	}, s3.WithPresignExpires(*aws.Duration(time.Minute * 10)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

	if sessionId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	keyName := h.getKeyForBlobInFlight(imageName, uploadId)

	partNumber := 0

	fullBytesRead := 0

	eg := errgroup.Group{}
	eg.SetLimit(10)

	// Read file in chunks
	for {
		buffer := make([]byte, h.Config.ChunkBufferLength)
		bytesRead, readErr := io.ReadFull(r.Body, buffer)

		if bytesRead > 0 {
			partNumber++
			eg.Go(func() error {
				_, err := h.asyncPartUpload(r.Context(), keyName, buffer, bytesRead, partNumber, sessionId)
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
			http.Error(w, "Failed to read the file", http.StatusInternalServerError)
			return
		}
	}

	if err := eg.Wait(); err != nil {
		http.Error(w, "Failed to upload part to S3", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Range", fmt.Sprintf("0-%d", fullBytesRead))
	w.Header().Set("Location", h.getUploadUrl(imageName, uploadId, sessionId))
	w.WriteHeader(http.StatusAccepted)
}

func (h *V2BlobsHandler) FinalizeBlobUploadSession(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	uploadId := r.Context().Value("uploadId").(string)
	sessionId := r.URL.Query().Get("session")
	blobDigest := r.URL.Query().Get("digest")
	log.Println(blobDigest)

	keyName := h.getKeyForBlobInFlight(imageName, uploadId)

	partsOutput, err := h.S3Client.ListParts(r.Context(), &s3.ListPartsInput{
		Bucket:   &h.Config.S3.BlobsBucketName,
		Key:      &keyName,
		UploadId: &sessionId,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var completedParts []types.CompletedPart

	for _, p := range partsOutput.Parts {
		completedParts = append(completedParts, types.CompletedPart{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
		})
	}

	_, err = h.S3Client.CompleteMultipartUpload(r.Context(), &s3.CompleteMultipartUploadInput{
		Bucket:   &h.Config.S3.BlobsBucketName,
		Key:      &keyName,
		UploadId: &sessionId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	headResp, err := h.S3Client.HeadObject(r.Context(), &s3.HeadObjectInput{
		Bucket: &h.Config.S3.BlobsBucketName,
		Key:    &keyName,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	objectSize := *headResp.ContentLength

	if objectSize > (5 * GiB) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	destKey := h.getKeyForBlob(imageName, blobDigest)
	copySource := fmt.Sprintf("%s/%s", h.Config.S3.BlobsBucketName, keyName)

	_, err = h.S3Client.CopyObject(r.Context(), &s3.CopyObjectInput{
		Bucket:     &h.Config.S3.BlobsBucketName,
		CopySource: &copySource,
		Key:        &destKey,
	})

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", h.getBlobDownloadUrl(imageName, blobDigest))
	w.WriteHeader(http.StatusCreated)
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

func (h *V2BlobsHandler) getKeyForBlobInFlight(imageName string, uploadId string) string {
	return fmt.Sprintf("%s/in-flight/%s.blob", imageName, uploadId)
}

func (h *V2BlobsHandler) getKeyForBlob(imageName string, digest string) string {
	return fmt.Sprintf("%s/blobs/%s/blob.data", imageName, getDigestAsNestedFolder(digest))
}

func (h *V2BlobsHandler) getUploadUrl(imageName string, uploadId string, s3UploadId string) string {
	return fmt.Sprintf("%s/v2/%s/blobs/uploads/%s?session=%s", h.Config.BaseURL, imageName, uploadId, s3UploadId)
}

func (h *V2BlobsHandler) getBlobDownloadUrl(imageName string, digest string) string {
	return fmt.Sprintf("%s/v2/%s/blobs/%s", h.Config.BaseURL, imageName, digest)
}

func getDigestAsNestedFolder(digest string) string {
	// Remove the "sha256:" prefix if it exists
	if strings.HasPrefix(digest, sha256Prefix) {
		digest = strings.TrimPrefix(digest, sha256Prefix)
	}

	// Split the digest into chunks of 2 characters
	var folders []string
	for i := 0; i < len(digest); i += 2 {
		if i+2 > len(digest) {
			// Handle any remainder (unlikely with SHA256 as it's 64 characters)
			folders = append(folders, digest[i:])
		} else {
			folders = append(folders, digest[i:i+2])
		}
	}

	// Join the folders with a "/" to simulate the S3 folder structure
	return strings.Join(folders, "/")
}

func isSHA256Digest(digest string) bool {
	return strings.HasPrefix(digest, sha256Prefix)
}
