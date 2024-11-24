package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
)

type V2ManifestsHandler struct {
	Config          *config.AppConfig
	S3Client        *s3.Client
	S3PresignClient *s3.PresignClient
}

func (h *V2ManifestsHandler) UploadManifest(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	reference := r.Context().Value("reference").(string)
	keyName := h.getKeyForManifestRef(imageName, reference)

	defer r.Body.Close()

	buffer, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	contentLength := int64(len(buffer))

	digest := reference
	if !isSHA256Digest(reference) {
		hash := sha256.New()
		hash.Write(buffer)
		hashDigest := hash.Sum(nil)
		digest = fmt.Sprintf("%x", hashDigest)
	}
	checksumBytes, err := hex.DecodeString(digest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	digestEnc := base64.StdEncoding.EncodeToString(checksumBytes)

	_, err = h.S3Client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket:         &h.Config.S3.BlobsBucketName,
		Key:            &keyName,
		Body:           strings.NewReader(string(buffer)),
		ContentLength:  &contentLength,
		ContentType:    aws.String("application/json"),
		ChecksumSHA256: &digestEnc,
	})

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Docker-Content-Digest", fmt.Sprintf("sha256:%s", digest))
	w.Header().Set("Location", h.getManifestDownloadUrl(imageName, reference))
	w.WriteHeader(http.StatusCreated)
}

func (h *V2ManifestsHandler) DownloadManifest(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	reference := r.Context().Value("reference").(string)

	keyName := h.getKeyForManifestRef(imageName, reference)

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

func (h *V2ManifestsHandler) HeadManifest(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	reference := r.Context().Value("reference").(string)

	output, err := h.S3Client.HeadObject(r.Context(), &s3.HeadObjectInput{
		Bucket:       &h.Config.S3.BlobsBucketName,
		Key:          aws.String(h.getKeyForManifestRef(imageName, reference)),
		ChecksumMode: types.ChecksumModeEnabled,
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

	checksumBytes, err := base64.StdEncoding.DecodeString(*output.ChecksumSHA256)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	digest := hex.EncodeToString(checksumBytes)

	w.Header().Set("Content-Length", strconv.Itoa(int(*output.ContentLength)))
	w.Header().Set("Docker-Content-Digest", digest)
	w.WriteHeader(http.StatusOK)
}

func (h *V2ManifestsHandler) getKeyForManifestRef(imageName string, reference string) string {
	ref := reference
	if isSHA256Digest(reference) {
		ref = getDigestAsNestedFolder(reference)
	}
	return fmt.Sprintf("%s/manifests/%s/manifest.json", imageName, ref)
}

func (h *V2ManifestsHandler) getManifestDownloadUrl(imageName string, reference string) string {
	return fmt.Sprintf("%s/v2/%s/manifests/%s", h.Config.BaseURL, imageName, reference)
}
