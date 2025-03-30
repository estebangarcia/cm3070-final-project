package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/go-chi/chi/v5"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

type OrganizationMiddleware struct {
	Config                 *config.AppConfig
	JwkCache               *jwk.Cache
	OrganizationRepository *repositories.OrganizationRepository
	RegistryRepository     *repositories.RegistryRepository
}

// This middleware validates that the authenticated user belongs to the organization
// is trying to access, if not it returns a 404
func (a *OrganizationMiddleware) ValidateOrg(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userSub := r.Context().Value("user_sub").(string)
		orgSlug := chi.URLParam(r, "organizationSlug")

		// Get the organization if the user belongs to it
		org, found, err := a.OrganizationRepository.GetForUserAndSlug(r.Context(), userSub, orgSlug)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		if !found {
			w.WriteHeader(404)
			return
		}

		// Store the organization in the context to be used in the handlers if needed
		ctx := context.WithValue(r.Context(), "organization", org)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// This middleware validates that the authenticated user belongs to the organization and has access to
// the registry specified in the URL
func (a *OrganizationMiddleware) ValidateOrgAndRegistry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userSubToken := r.Context().Value("user_sub").(string)
		orgSlug := chi.URLParam(r, "organizationSlug")
		registrySlug := chi.URLParam(r, "registrySlug")

		// if the user is an admin then we have access to everything
		if userSubToken == a.Config.AdminUser.Sub {
			userSubToken = ""
		}

		// Get the registry if the user has access to it
		registry, found, err := a.RegistryRepository.GetForOrgAndUser(r.Context(), userSubToken, orgSlug, registrySlug)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		if !found {
			w.WriteHeader(404)
			return
		}

		// Store the registry and organization in the context to be used in handlers if needed
		ctx := context.WithValue(r.Context(), "registry", registry)
		ctx = context.WithValue(ctx, "organization", registry.Edges.Organization)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
