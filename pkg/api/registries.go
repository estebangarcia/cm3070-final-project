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

// Get a specific registry by name, we can fetch this from the context as it
// is already parsed and queried by the middleware
func (rh *RegistriesHandler) GetRegistry(w http.ResponseWriter, r *http.Request) {
	registry := r.Context().Value("registry").(*ent.Registry)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(registry)
}

// Get all the registries for the specified organization
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

// Create a registry by name
func (rh *RegistriesHandler) CreateRegistries(w http.ResponseWriter, r *http.Request) {
	org := r.Context().Value("organization").(*ent.Organization)

	// Validate creation request
	createRegistryRequest, err := requests.BindRequest[requests.CreateRegistryRequest](r)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		return
	}

	// Create registry in the database
	registry, err := rh.RegistryRepository.CreateRegistry(r.Context(), createRegistryRequest.Name, org.ID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registry)
}
