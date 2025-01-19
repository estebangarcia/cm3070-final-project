package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/oci_models"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
)

type V2ReferrersHandler struct {
	Config               *config.AppConfig
	S3Client             *s3.Client
	S3PresignClient      *s3.PresignClient
	RepositoryRepository *repositories.RepositoryRepository
	ManifestRepository   *repositories.ManifestRepository
}

func (h *V2ReferrersHandler) GetReferrersForDigest(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("repositoryName").(string)
	reference := r.Context().Value("reference").(string)
	registry := r.Context().Value("registry").(*ent.Registry)
	artifactType := r.URL.Query().Get("artifactType")

	repo, found, err := h.RepositoryRepository.GetForRegistryByName(r.Context(), registry.ID, imageName)
	if err != nil {
		responses.OCIInternalServerError(w)
		return
	}

	if !found {
		responses.OCIRepositoryUnknown(w, imageName, true)
		return
	}

	referrers, err := h.ManifestRepository.GetManifestReferrers(r.Context(), reference, artifactType, repo)
	if err != nil {
		responses.OCIInternalServerError(w)
		return
	}

	manifests := []oci_models.OCIV1Manifest{}

	for _, ref := range referrers {
		output, err := h.S3Client.GetObject(r.Context(), &s3.GetObjectInput{
			Bucket: &h.Config.S3.BlobsBucketName,
			Key:    &ref.S3Path,
		})
		if err != nil {
			responses.OCIInternalServerError(w)
			return
		}

		manifest, err := decodeBytes[oci_models.OCIV1Manifest](output.Body)
		if err != nil {
			responses.OCIInternalServerError(w)
			return
		}

		manifest.Digest = &ref.Digest
		manifests = append(manifests, *manifest)
	}

	index := oci_models.NewOCIV1ManifestIndex(manifests)

	if artifactType != "" {
		w.Header().Set("OCI-Filters-Applied", "artifactType")
	}
	w.Header().Set("Content-Type", index.MediaType)
	json.NewEncoder(w).Encode(index)
}

func decodeBytes[T any](r io.Reader) (*T, error) {
	var requestData T

	err := json.NewDecoder(r).Decode(&requestData)
	return &requestData, err
}
