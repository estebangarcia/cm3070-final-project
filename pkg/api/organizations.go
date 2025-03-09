package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/requests"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
	"github.com/go-chi/chi/v5"
)

type OrganizationsHandler struct {
	Config                 *config.AppConfig
	OrganizationRepository *repositories.OrganizationRepository
	UserRepository         *repositories.UserRepository
	RegistryRepository     *repositories.RegistryRepository
	RepositoryRepository   *repositories.RepositoryRepository
	ManifestRepository     *repositories.ManifestRepository
}

func (oh *OrganizationsHandler) GetOrganizationsForUser(w http.ResponseWriter, r *http.Request) {
	userSub := r.Context().Value("user_sub").(string)

	orgs, err := oh.OrganizationRepository.GetForUser(r.Context(), userSub)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orgs)
}

func (oh *OrganizationsHandler) GetOrganizationsBySlugForUser(w http.ResponseWriter, r *http.Request) {
	userSub := r.Context().Value("user_sub").(string)
	orgSlug := chi.URLParam(r, "organizationSlug")

	org, found, err := oh.OrganizationRepository.GetForUserAndSlug(r.Context(), userSub, orgSlug)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if !found {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(org)
}

func (oh *OrganizationsHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	userSub := r.Context().Value("user_sub").(string)

	createOrgRequest, err := requests.BindRequest[requests.CreateOrganizationRequest](r)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		return
	}

	user, err := oh.UserRepository.GetUserBySub(r.Context(), userSub)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	org, err := oh.OrganizationRepository.CreateOrganizationWithOwner(r.Context(), user, createOrgRequest.Name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(org)
}

func (oh *OrganizationsHandler) GetOrganizationStats(w http.ResponseWriter, r *http.Request) {
	organization := r.Context().Value("organization").(*ent.Organization)

	storageUsed, err := oh.ManifestRepository.GetStorageUsedInBytesForOrganization(r.Context(), organization)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	registryCount, err := oh.RegistryRepository.GetCountForOrg(r.Context(), organization)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	repositoryCount, err := oh.RepositoryRepository.GetCountForOrg(r.Context(), organization)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	artifactCount, err := oh.ManifestRepository.GetCountForOrg(r.Context(), organization)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	artifactWithVulnerabilitiesCount, err := oh.ManifestRepository.GetCountWithVulnerabilitiesForOrg(r.Context(), organization)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	response := &responses.OrganizationStatsResponse{
		RegistryCount:            registryCount,
		RepositoryCount:          repositoryCount,
		ArtifactsCount:           artifactCount,
		StorageUsed:              storageUsed,
		VulnerableArtifactsCount: artifactWithVulnerabilitiesCount,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
