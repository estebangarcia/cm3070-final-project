package api

import (
	"context"
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
	"github.com/estebangarcia/cm3070-final-project/pkg/oci_models"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/requests"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
)

type V2ManifestsHandler struct {
	Config                *config.AppConfig
	S3Client              *s3.Client
	S3PresignClient       *s3.PresignClient
	RepositoryRepository  *repositories.RepositoryRepository
	ManifestRepository    *repositories.ManifestRepository
	ManifestTagRepository *repositories.ManifestTagRepository
}

func (h *V2ManifestsHandler) UploadManifest(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("repositoryName").(string)
	reference := r.Context().Value("reference").(string)
	org := r.Context().Value("organization").(*ent.Organization)
	registry := r.Context().Value("registry").(*ent.Registry)
	manifestContentType := r.Header.Get("Content-Type")

	repo, err := h.RepositoryRepository.GetOrCreateRepository(r.Context(), registry.ID, imageName)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}
	defer r.Body.Close()

	buffer, err := io.ReadAll(r.Body)
	if err != nil {
		responses.OCIInternalServerError(w)
		return
	}

	manifestRequest, err := requests.BindRequestFromBytes[oci_models.OCIV1Manifest](buffer)
	if err != nil {
		responses.OCIUnprocessableEntity(w, "the manifest is bad")
		return
	}

	contentLength := int64(len(buffer))
	digest, checksumBytes, err := h.getDigestFromReferenceOrBody(reference, buffer)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	digestWithPrefix := fmt.Sprintf("sha256:%s", digest)

	keyName := h.getKeyForManifestRef(manifestContentType, org.Slug, registry.Slug, imageName, digestWithPrefix)

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

	var subjectManifest *ent.Manifest = nil
	var subjectManifestFound bool

	// This manifest is trying to refer to another one, check if it exists if it doesn't then create a placeholder for a future manifest
	if manifestRequest.Subject != nil {
		subjectManifest, subjectManifestFound, err = h.ManifestRepository.GetManifestByReferenceAndMediaType(r.Context(), manifestRequest.Subject.Digest, manifestRequest.Subject.MediaType, repo)
		if err != nil {
			log.Println(err)
			responses.OCIInternalServerError(w)
			return
		}

		if !subjectManifestFound {
			subjectManifest = &ent.Manifest{}
			subjectManifest.Digest = manifestRequest.Subject.Digest
			subjectManifest.MediaType = manifestRequest.Subject.MediaType
		}
	}

	artifactType := manifestRequest.ArtifactType
	if artifactType == nil {
		artifactType = &manifestRequest.Config.MediaType
	}

	var layers []*ent.ManifestLayer
	for _, layer := range manifestRequest.Layers {
		layers = append(layers, &ent.ManifestLayer{
			MediaType:   layer.MediaType,
			Digest:      layer.Digest,
			Annotations: layer.Annotations,
			Size:        int32(layer.Size),
		})
	}

	m, err := h.ManifestRepository.UpsertManifestWithSubjectAndTag(r.Context(), layers, reference, digestWithPrefix, manifestRequest.MediaType, artifactType, keyName, subjectManifest, repo)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	w.Header().Set("Docker-Content-Digest", m.Digest)
	w.Header().Set("Location", h.getManifestDownloadUrl(org.Slug, registry.Slug, imageName, reference))
	if subjectManifest != nil {
		w.Header().Set("OCI-Subject", subjectManifest.Digest)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *V2ManifestsHandler) DownloadManifest(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("repositoryName").(string)
	reference := r.Context().Value("reference").(string)
	registry := r.Context().Value("registry").(*ent.Registry)
	acceptedTypes := h.getAcceptedTypes(r)

	repo, err := h.RepositoryRepository.GetOrCreateRepository(r.Context(), registry.ID, imageName)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	manifests, err := h.ManifestRepository.GetManifestsByReferenceAndMediaType(r.Context(), reference, acceptedTypes, repo)
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
	imageName := r.Context().Value("repositoryName").(string)
	reference := r.Context().Value("reference").(string)
	registry := r.Context().Value("registry").(*ent.Registry)
	acceptedTypes := h.getAcceptedTypes(r)

	repo, err := h.RepositoryRepository.GetOrCreateRepository(r.Context(), registry.ID, imageName)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	manifests, err := h.ManifestRepository.GetManifestsByReferenceAndMediaType(r.Context(), reference, acceptedTypes, repo)
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

func (h *V2ManifestsHandler) DeleteManifestOrTag(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("repositoryName").(string)
	reference := r.Context().Value("reference").(string)
	registry := r.Context().Value("registry").(*ent.Registry)

	repo, found, err := h.RepositoryRepository.GetForRegistryByName(r.Context(), registry.ID, imageName)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	if !found {
		responses.OCIRepositoryUnknown(w, imageName, false)
		return
	}

	if !helpers.IsSHA256Digest(reference) {
		h.deleteTag(r.Context(), w, repo, reference)
		return
	}

	h.deleteManifestByDigest(r.Context(), w, repo, reference)
}

func (h *V2ManifestsHandler) deleteTag(ctx context.Context, w http.ResponseWriter, repo *ent.Repository, tagRef string) {
	tag, found, err := h.ManifestTagRepository.GetTagByName(ctx, repo, tagRef)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	if !found {
		responses.OCITagUnknown(w, repo.Name, tagRef)
		return
	}

	if err := h.ManifestTagRepository.DeleteTag(ctx, tag); err != nil {
		responses.OCIInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *V2ManifestsHandler) deleteManifestByDigest(ctx context.Context, w http.ResponseWriter, repo *ent.Repository, digest string) {
	manifest, found, err := h.ManifestRepository.GetManifestByReference(ctx, digest, repo, false)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	if !found {
		responses.OCIManifestUnknown(w, digest)
		return
	}

	if err := h.ManifestRepository.DeleteManifest(ctx, manifest); err != nil {
		responses.OCIInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *V2ManifestsHandler) getKeyForManifestRef(contentType string, orgSlug string, registrySlug string, imageName string, digest string) string {
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

	return fmt.Sprintf("%s/manifests/%s/%s/%s%s/manifest.json", orgSlug, registrySlug, imageName, digest, contentTypeSubFolder)
}

func (h *V2ManifestsHandler) getManifestDownloadUrl(orgSlug string, registrySlug string, imageName string, reference string) string {
	return fmt.Sprintf("%s/v2/%s/%s/%s/manifests/%s", h.Config.GetBaseUrl(), orgSlug, registrySlug, imageName, reference)
}

func (h *V2ManifestsHandler) getDigestFromReferenceOrBody(reference string, body []byte) (string, []byte, error) {
	digest := helpers.TrimDigest(reference)
	if !helpers.IsSHA256Digest(reference) {
		hash := sha256.New()
		hash.Write(body)
		hashDigest := hash.Sum(nil)
		digest = fmt.Sprintf("%x", hashDigest)
	}

	checksumBytes, err := hex.DecodeString(digest)

	return digest, checksumBytes, err
}

func (h *V2ManifestsHandler) getAcceptedTypes(r *http.Request) []string {
	acceptedTypes := r.Header["Accept"]
	if len(acceptedTypes) == 1 && strings.Contains(acceptedTypes[0], ",") {
		acceptedTypesComma := acceptedTypes[0]
		acceptedTypesComma = strings.ReplaceAll(acceptedTypesComma, " ", "")
		acceptedTypes = strings.Split(acceptedTypesComma, ",")
	}

	for i, acceptedType := range acceptedTypes {
		if acceptedType == "*/*" {
			return append(acceptedTypes[:i], acceptedTypes[i+1:]...)
		}
	}

	return acceptedTypes
}
