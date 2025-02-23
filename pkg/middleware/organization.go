package middleware

import (
	"context"
	"fmt"
	"log"
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

func (a *OrganizationMiddleware) ValidateOrg(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userSub := r.Context().Value("user_sub").(string)
		orgSlug := chi.URLParam(r, "organizationSlug")

		org, found, err := a.OrganizationRepository.GetForUserAndSlug(r.Context(), userSub, orgSlug)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		if !found {
			w.WriteHeader(404)
			return
		}

		ctx := context.WithValue(r.Context(), "organization", org)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *OrganizationMiddleware) ValidateOrgAndRegistry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userSubToken := r.Context().Value("user_sub").(string)
		orgSlug := chi.URLParam(r, "organizationSlug")
		registrySlug := chi.URLParam(r, "registrySlug")

		if userSubToken == a.Config.AdminUser.Sub {
			userSubToken = ""
		}

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

		ctx := context.WithValue(r.Context(), "registry", registry)
		ctx = context.WithValue(ctx, "organization", registry.Edges.Organization)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/*func (a *OrganizationMiddleware) ValidateOrgRegistryRepository(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userSub := r.Context().Value("user_sub").(string)
		orgSlug := chi.URLParam(r, "organizationSlug")
		registrySlug := chi.URLParam(r, "registrySlug")
		repositoryName := chi.URLParam(r, "repositoryName")

		registry, found, err := a.RegistryRepository.GetForOrgAndUser(r.Context(), userSub, orgSlug, registrySlug)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		if !found {
			w.WriteHeader(404)
			return
		}

		ctx := context.WithValue(r.Context(), "registry", registry)
		ctx = context.WithValue(ctx, "organization", registry.Edges.Organization)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
*/
