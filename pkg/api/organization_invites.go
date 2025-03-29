package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/go-chi/chi/v5"
)

type OrganizationInvitesHandler struct {
	Config                       *config.AppConfig
	OrganizationInviteRepository *repositories.OrganizationInviteRepository
}

// Get all the invites for a user
func (oh *OrganizationInvitesHandler) GetInvitesForUser(w http.ResponseWriter, r *http.Request) {
	userSub := r.Context().Value("user_sub").(string)

	invites, err := oh.OrganizationInviteRepository.GetInvitesForUser(r.Context(), userSub)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invites)
}

// Accept an invitation
func (oh *OrganizationInvitesHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	userSub := r.Context().Value("user_sub").(string)
	inviteId := chi.URLParam(r, "inviteId")

	hasInvite, err := oh.OrganizationInviteRepository.HasInviteWithID(r.Context(), inviteId, userSub)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if !hasInvite {
		w.WriteHeader(404)
		return
	}

	if err := oh.OrganizationInviteRepository.AcceptInvite(r.Context(), inviteId, userSub); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Reject an invitation
func (oh *OrganizationInvitesHandler) RejectInvite(w http.ResponseWriter, r *http.Request) {
	userSub := r.Context().Value("user_sub").(string)
	inviteId := chi.URLParam(r, "inviteId")

	hasInvite, err := oh.OrganizationInviteRepository.HasInviteWithID(r.Context(), inviteId, userSub)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if !hasInvite {
		w.WriteHeader(404)
		return
	}

	if err = oh.OrganizationInviteRepository.RejectInvite(r.Context(), inviteId, userSub); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
