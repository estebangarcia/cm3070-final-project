package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
)

type RepositoriesHandler struct {
	Config               *config.AppConfig
	RepositoryRepository *repositories.RepositoryRepository
}

func (rh *RepositoriesHandler) GetRepository(w http.ResponseWriter, r *http.Request) {
	registry := r.Context().Value("registry").(*ent.Registry)
	repositoryName := r.Context().Value("repositoryName").(string)

	repo, found, err := rh.RepositoryRepository.GetForRegistryByName(r.Context(), registry.ID, repositoryName)
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
	json.NewEncoder(w).Encode(repo)
}

func (rh *RepositoriesHandler) GetRepositories(w http.ResponseWriter, r *http.Request) {
	registry := r.Context().Value("registry").(*ent.Registry)

	repositories, err := rh.RepositoryRepository.GetAllForRegistry(r.Context(), registry.ID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repositories)
}
