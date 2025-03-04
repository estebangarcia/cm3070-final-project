package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/requests"
)

type RegistriesHandler struct {
	Config             *config.AppConfig
	RegistryRepository *repositories.RegistryRepository
}

func (rh *RegistriesHandler) GetRegistry(w http.ResponseWriter, r *http.Request) {
	registry := r.Context().Value("registry").(*ent.Registry)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(registry)
}

func (rh *RegistriesHandler) GetRegistries(w http.ResponseWriter, r *http.Request) {
	org := r.Context().Value("organization").(*ent.Organization)

	registries, err := rh.RegistryRepository.GetForOrg(r.Context(), org.Slug)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(registries)
}

func (rh *RegistriesHandler) CreateRegistries(w http.ResponseWriter, r *http.Request) {
	org := r.Context().Value("organization").(*ent.Organization)

	createRegistryRequest, err := requests.BindRequest[requests.CreateRegistryRequest](r)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		return
	}

	registry, err := rh.RegistryRepository.CreateRegistry(r.Context(), createRegistryRequest.Name, org.ID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registry)
}
