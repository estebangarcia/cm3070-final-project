package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
)

type VulnerabilitiesHandlers struct {
	Config               *config.AppConfig
	RepositoryRepository *repositories.RepositoryRepository
	ManifestRepository   *repositories.ManifestRepository
}

// Get all the vulnerabilities for a specific artifact by its manifest digest
// it returns both misconfigurations and vulnerabilities
func (vh *VulnerabilitiesHandlers) GetVulnerabilitiesForArtifact(w http.ResponseWriter, r *http.Request) {
	registry := r.Context().Value("registry").(*ent.Registry)
	repositoryName := r.Context().Value("repositoryName").(string)
	manifestDigest := r.Context().Value("digest").(string)

	repo, found, err := vh.RepositoryRepository.GetForRegistryByName(r.Context(), registry.ID, repositoryName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	vulnerabilities, err := vh.ManifestRepository.GetManifestVulnerabilitiesByReference(r.Context(), manifestDigest, repo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	misconfigurations, err := vh.ManifestRepository.GetManifestMisconfigurationsByReference(r.Context(), manifestDigest, repo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.VulnerabilitiesResponse{
		Vulnerabilities:   vulnerabilities,
		Misconfigurations: misconfigurations,
	})
}
