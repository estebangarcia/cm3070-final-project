package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
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
