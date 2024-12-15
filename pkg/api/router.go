package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/helpers"
	"github.com/estebangarcia/cm3070-final-project/pkg/middleware"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	engine *chi.Mux
}

func NewRouter(ctx context.Context, cfg config.AppConfig, dbClient *ent.Client) (*Router, error) {
	r := chi.NewRouter()

	apiCustomMux := NewCustomMux()
	v2CustomMux := NewCustomMux()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	/* DB Repositories */
	blobChunkRepository := repositories.NewBlobChunkRepository(dbClient)
	manifestRepository := repositories.NewManifestRepository(dbClient)
	repositoryRepository := repositories.NewRepositoryRepository(dbClient)
	organizationsRepository := repositories.NewOrganizationRepository(dbClient)
	registryRepository := repositories.NewRegistryRepository(dbClient)

	/* AWS Helpers */
	s3Client := helpers.GetS3Client(ctx, cfg)
	s3Presigner := helpers.GetS3PresignClient(s3Client)
	jwkCache, err := helpers.InitJWKCache(ctx, &cfg)
	if err != nil {
		return nil, err
	}

	/* Middlewares */
	jwtAuthMiddleware := middleware.JWTAuthMiddleware{
		Config:   &cfg,
		JwkCache: jwkCache,
	}
	organizationsHandlers := OrganizationsHandler{
		Config:                 &cfg,
		OrganizationRepository: organizationsRepository,
	}

	orgMiddleware := middleware.OrganizationMiddleware{
		Config:                 &cfg,
		OrganizationRepository: organizationsRepository,
		RegistryRepository:     registryRepository,
	}
	extractBasicAuthMiddleware := middleware.ExtractBasicCredentialsMiddleware{}

	/* HTTP Handlers */
	healthHandler := HealthHandler{}

	v2LoginHandler := V2LoginHandler{
		Config:        &cfg,
		CognitoClient: helpers.GetCognitoClient(ctx, cfg),
	}

	v2PingHandler := V2PingHandler{}

	v2BlobsHandler := V2BlobsHandler{
		Config:              &cfg,
		S3Client:            s3Client,
		S3PresignClient:     s3Presigner,
		BlobChunkRepository: blobChunkRepository,
	}

	v2ManifestsHandlers := V2ManifestsHandler{
		Config:               &cfg,
		S3Client:             s3Client,
		S3PresignClient:      s3Presigner,
		ManifestRepository:   manifestRepository,
		RepositoryRepository: repositoryRepository,
	}

	registriesHandler := RegistriesHandler{
		Config:             &cfg,
		RegistryRepository: registryRepository,
	}

	repositoriesHandler := RepositoriesHandler{
		Config:               &cfg,
		RepositoryRepository: repositoryRepository,
	}

	r.Get("/api/v1/health", healthHandler.GetHealth)

	r.Route("/api/v1", func(authenticatedApiV1 chi.Router) {
		authenticatedApiV1.Use(jwtAuthMiddleware.Validate)
		authenticatedApiV1.Get("/organizations", organizationsHandlers.GetOrganizationsForUser)

		authenticatedApiV1.Route("/organizations/{organizationSlug:[a-z0-9-]+}", func(orgScopedRoutes chi.Router) {
			orgScopedRoutes.Use(orgMiddleware.ValidateOrg)
			orgScopedRoutes.Get("/", organizationsHandlers.GetOrganizationsBySlugForUser)
			orgScopedRoutes.Get("/registries", registriesHandler.GetRegistries)
			orgScopedRoutes.Post("/registries", registriesHandler.CreateRegistries)
		})

		authenticatedApiV1.Route("/organizations/{organizationSlug:[a-z0-9-]+}/registries/{registrySlug:[a-z0-9-]+}", func(registryScopedRoutes chi.Router) {
			registryScopedRoutes.Use(orgMiddleware.ValidateOrgAndRegistry)
			registryScopedRoutes.Get("/", registriesHandler.GetRegistry)
		})

		authenticatedApiV1.Route("/organizations/{organizationSlug:[a-z0-9-]+}/registries/{registrySlug:[a-z0-9-]+}/repositories", func(repositoryScopedRoutes chi.Router) {
			repositoryScopedRoutes.Use(orgMiddleware.ValidateOrgAndRegistry)
			repositoryScopedRoutes.Get("/", repositoriesHandler.GetRepositories)

			apiCustomMux.Get(getRepositoryRegexRoute(), repositoriesHandler.GetRepository)
			repositoryScopedRoutes.HandleFunc("/*", apiCustomMux.Handle)
		})
	})

	r.With(extractBasicAuthMiddleware.Validate).Get("/v2/login", v2LoginHandler.Login)
	r.Route("/v2", func(authenticatedOciV2 chi.Router) {
		authenticatedOciV2.Use(jwtAuthMiddleware.Validate)
		authenticatedOciV2.Get("/", v2PingHandler.Ping)
		authenticatedOciV2.Route("/{organizationSlug:[a-z-]+}/{registrySlug:[a-z-]+}", func(registryScopedOCIRoutes chi.Router) {
			registryScopedOCIRoutes.Use(orgMiddleware.ValidateOrgAndRegistry)

			v2CustomMux.Post(getRepositoryRegexRoute()+`\/blobs\/uploads`, v2BlobsHandler.InitiateUploadSession)
			v2CustomMux.Patch(getRepositoryRegexRoute()+`\/blobs\/uploads\/(?P<uploadId>[\w]+)`, v2BlobsHandler.UploadBlob)
			v2CustomMux.Put(getRepositoryRegexRoute()+`\/blobs\/uploads\/(?P<uploadId>[\w]+)`, v2BlobsHandler.FinalizeBlobUploadSession)
			v2CustomMux.Get(getRepositoryRegexRoute()+`\/blobs\/uploads\/(?P<uploadId>[\w]+)`, v2BlobsHandler.GetBlobUploadSession)
			v2CustomMux.Get(getRepositoryRegexRoute()+`\/blobs\/(?P<digest>[\/\w:]+)`, v2BlobsHandler.DownloadBlob)
			v2CustomMux.Head(getRepositoryRegexRoute()+`\/blobs\/(?P<digest>[\/\w:]+)`, v2BlobsHandler.HeadBlob)

			v2CustomMux.Put(getRepositoryRegexRoute()+`\/manifests\/(?P<reference>[\w:._-]+)`, v2ManifestsHandlers.UploadManifest)
			v2CustomMux.Get(getRepositoryRegexRoute()+`\/manifests\/(?P<reference>[\w:._-]+)`, v2ManifestsHandlers.DownloadManifest)
			v2CustomMux.Head(getRepositoryRegexRoute()+`\/manifests\/(?P<reference>[\w:._-]+)`, v2ManifestsHandlers.HeadManifest)

			registryScopedOCIRoutes.HandleFunc("/*", v2CustomMux.Handle)
		})
	})

	return &Router{
		engine: r,
	}, nil
}

func (r Router) Run(ctx context.Context, portBinding string) error {
	srv := &http.Server{
		Addr:              portBinding,
		Handler:           r.engine,
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: 3 * time.Second,
	}

	go r.listen(srv)
	<-ctx.Done()
	log.Println("shutting down server")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(timeoutCtx); err != nil {
		return err
	}

	return nil
}

func (r Router) listen(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error trying to listen: %s\n", err)
	}
}

func getRepositoryRegexRoute() string {
	return fmt.Sprintf(`(?P<repositoryName>%s)`, `[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*(\/[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*)*`)
}
