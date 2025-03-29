package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
)

type ArtifactsHandler struct {
	Config               *config.AppConfig
	RepositoryRepository *repositories.RepositoryRepository
	ManifestRepository   *repositories.ManifestRepository
}

// Get all artifacts for a specified repository
func (ah *ArtifactsHandler) GetArtifactsForRepository(w http.ResponseWriter, r *http.Request) {
	registry := r.Context().Value("registry").(*ent.Registry)
	repositoryName := r.Context().Value("repositoryName").(string)

	repo, found, err := ah.RepositoryRepository.GetForRegistryByName(r.Context(), registry.ID, repositoryName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// returns all manifests with their associated tags
	manifests, err := ah.ManifestRepository.GetAllWithTags(r.Context(), repo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(manifests)
}

// Get a specific artifact by its manifest digest
func (ah *ArtifactsHandler) GetArtifactByDigest(w http.ResponseWriter, r *http.Request) {
	registry := r.Context().Value("registry").(*ent.Registry)
	repositoryName := r.Context().Value("repositoryName").(string)
	manifestDigest := r.Context().Value("digest").(string)

	repo, found, err := ah.RepositoryRepository.GetForRegistryByName(r.Context(), registry.ID, repositoryName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	manifest, found, err := ah.ManifestRepository.GetManifestByReference(r.Context(), manifestDigest, repo, true)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(manifest)
}

// Get all artifacts for an organization
func (ah *ArtifactsHandler) GetArtifactsForOrg(w http.ResponseWriter, r *http.Request) {
	organization := r.Context().Value("organization").(*ent.Organization)

	manifests, err := ah.ManifestRepository.GetAllForOrgWithTags(r.Context(), organization)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(manifests)
}
