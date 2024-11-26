package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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
	"github.com/estebangarcia/cm3070-final-project/pkg/helpers"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
)

type V2ManifestsHandler struct {
	Config               *config.AppConfig
	S3Client             *s3.Client
	S3PresignClient      *s3.PresignClient
	RepositoryRepository *repositories.RepositoryRepository
	ManifestRepository   *repositories.ManifestRepository
}

func (h *V2ManifestsHandler) UploadManifest(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	reference := r.Context().Value("reference").(string)
	manifestContentType := r.Header.Get("Content-Type")

	repo, err := h.RepositoryRepository.GetOrCreateRepository(r.Context(), imageName)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	keyName := h.getKeyForManifestRef(manifestContentType, imageName, reference)

	defer r.Body.Close()

	buffer, err := io.ReadAll(r.Body)
	if err != nil {
		responses.OCIInternalServerError(w)
		return
	}

	contentLength := int64(len(buffer))

	digest, checksumBytes, err := h.getDigestFromReferenceOrBody(reference, buffer)
	if err != nil {
		responses.OCIInternalServerError(w)
		return
	}

	digestEnc := base64.StdEncoding.EncodeToString(checksumBytes)

	_, err = h.S3Client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket:         &h.Config.S3.BlobsBucketName,
		Key:            &keyName,
		Body:           strings.NewReader(string(buffer)),
		ContentLength:  &contentLength,
		ContentType:    aws.String(manifestContentType),
		ChecksumSHA256: &digestEnc,
	})

	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	digestWithPrefix := fmt.Sprintf("sha256:%s", digest)

	m, err := h.ManifestRepository.CreateManifestAndUpsertTag(r.Context(), reference, digestWithPrefix, manifestContentType, keyName, repo)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	w.Header().Set("Docker-Content-Digest", m.Digest)
	w.Header().Set("Location", h.getManifestDownloadUrl(imageName, reference))
	w.WriteHeader(http.StatusCreated)
}

func (h *V2ManifestsHandler) DownloadManifest(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	reference := r.Context().Value("reference").(string)
	acceptedTypes := r.Header.Values("Accept")

	manifests, err := h.ManifestRepository.GetManifestsByReferenceAndMediaType(r.Context(), reference, acceptedTypes, imageName)
	if err != nil {
		responses.OCIInternalServerError(w)
		log.Println(err)
		return
	}
	if len(manifests) == 0 {
		responses.OCIManifestUnknown(w, reference)
		return
	}
	manifest := manifests[0]

	req, err := h.S3PresignClient.PresignGetObject(r.Context(), &s3.GetObjectInput{
		Bucket: &h.Config.S3.BlobsBucketName,
		Key:    &manifest.S3Path,
	}, s3.WithPresignExpires(*aws.Duration(time.Minute * 10)))
	if err != nil {
		responses.OCIInternalServerError(w)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", manifest.MediaType)
	w.Header().Set("Location", req.URL)
	w.Header().Set("Docker-Content-Digest", manifest.Digest)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *V2ManifestsHandler) HeadManifest(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName").(string)
	reference := r.Context().Value("reference").(string)
	acceptedTypes := r.Header.Values("Accept")

	manifests, err := h.ManifestRepository.GetManifestsByReferenceAndMediaType(r.Context(), reference, acceptedTypes, imageName)
	if err != nil {
		responses.OCIInternalServerError(w)
		log.Println(err)
		return
	}
	if len(manifests) == 0 {
		responses.OCIManifestUnknown(w, reference)
		return
	}
	manifest := manifests[0]

	output, err := h.S3Client.HeadObject(r.Context(), &s3.HeadObjectInput{
		Bucket:       &h.Config.S3.BlobsBucketName,
		Key:          &manifest.S3Path,
		ChecksumMode: types.ChecksumModeEnabled,
	})
	if err != nil {
		responses.OCIInternalServerError(w)
		log.Println(err)
	}

	w.Header().Set("Content-Type", manifest.MediaType)
	w.Header().Set("Content-Length", strconv.Itoa(int(*output.ContentLength)))
	w.Header().Set("Docker-Content-Digest", manifest.Digest)
	w.WriteHeader(http.StatusOK)
}

func (h *V2ManifestsHandler) getKeyForManifestRef(contentType string, imageName string, reference string) string {
	ref := reference
	if helpers.IsSHA256Digest(reference) {
		ref = helpers.GetDigestAsNestedFolder(reference)
	}

	contentTypeSubFolder := ""

	if contentType != "" {
		appMediaType := strings.Split(contentType, "/")
		if helpers.IsVendorSpecificContentType(contentType) {
			vendorContentType := appMediaType[1]
			vendorContentType = strings.ReplaceAll(vendorContentType, "+", ".")
			vndSplit := strings.Split(vendorContentType, ".")
			vndComplete := append([]string{appMediaType[0]}, vndSplit...)

			contentTypeSubFolder = "/" + strings.Join(vndComplete, "/")
		} else {
			contentTypeSubFolder = "/" + contentType
		}
	}

	return fmt.Sprintf("%s/manifests/%s%s/manifest.json", imageName, ref, contentTypeSubFolder)
}

func (h *V2ManifestsHandler) getManifestDownloadUrl(imageName string, reference string) string {
	return fmt.Sprintf("%s/v2/%s/manifests/%s", h.Config.BaseURL, imageName, reference)
}

func (h *V2ManifestsHandler) getDigestFromReferenceOrBody(reference string, body []byte) (string, []byte, error) {
	digest := reference
	if !helpers.IsSHA256Digest(reference) {
		hash := sha256.New()
		hash.Write(body)
		hashDigest := hash.Sum(nil)
		digest = fmt.Sprintf("%x", hashDigest)
	}
	checksumBytes, err := hex.DecodeString(digest)

	return digest, checksumBytes, err
}
