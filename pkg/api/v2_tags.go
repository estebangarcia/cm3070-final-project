package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
)

type V2TagsHandler struct {
	Config                *config.AppConfig
	RepositoryRepository  *repositories.RepositoryRepository
	ManifestTagRepository *repositories.ManifestTagRepository
}

func (h *V2TagsHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("repositoryName").(string)
	registry := r.Context().Value("registry").(*ent.Registry)
	tagLimitNumber := r.URL.Query().Get("n")
	lastTag := r.URL.Query().Get("last")

	n := 10
	var err error

	if tagLimitNumber != "" {
		n, err = strconv.Atoi(tagLimitNumber)
		if err != nil {
			responses.OCIUnprocessableEntity(w, "'n' must be a number")
			return
		}
	}

	repo, exists, err := h.RepositoryRepository.GetForRegistryByName(r.Context(), registry.ID, imageName)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	if !exists {
		responses.OCIRepositoryUnknown(w, imageName)
		return
	}

	if lastTag != "" {
		_, exists, err = h.ManifestTagRepository.GetTagByName(r.Context(), repo, lastTag)
		if err != nil {
			log.Println(err)
			responses.OCIInternalServerError(w)
			return
		}

		if !exists {
			responses.OCITagUnknown(w, imageName, lastTag)
			return
		}
	}

	tags, err := h.ManifestTagRepository.ListTagsForRepository(r.Context(), repo, n, lastTag)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	t := []string{}

	for i, tag := range tags {
		if (i + 1) <= n {
			t = append(t, tag.Tag)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.TagsListResponse{
		Name: imageName,
		Tags: t,
	})
}
