package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/go-chi/chi/v5"
)

type OrganizationsHandler struct {
	Config                 *config.AppConfig
	OrganizationRepository *repositories.OrganizationRepository
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
