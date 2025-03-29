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

	// Create the custom multiplexor router
	apiCustomMux := NewCustomMux()
	v2CustomMux := NewCustomMux()

	// Add the chi Recoverer middleware, this recovers the server from crashing if there was a panic
	r.Use(chiMiddleware.Recoverer)
	// Add the chi strip slashes middleware, this removes the last slash of the URL
	r.Use(chiMiddleware.StripSlashes)

	/* Initialize all DB Repositories */
	blobChunkRepository := repositories.NewBlobChunkRepository()
	manifestRepository := repositories.NewManifestRepository()
	manifestTagRepository := repositories.NewManifestTagRepository()
	repositoryRepository := repositories.NewRepositoryRepository()
	organizationsRepository := repositories.NewOrganizationRepository()
	organizationInviteRepository := repositories.NewOrganizationInviteRepository()
	registryRepository := repositories.NewRegistryRepository()
	userRepository := repositories.NewUserRepository()

	/* Initialize the AWS SDK Clients needed  */
	s3Client := helpers.GetS3Client(ctx, cfg)
	s3Presigner := helpers.GetS3PresignClient(s3Client)
	sesClient := helpers.GetSESClient(ctx, cfg)
	jwkCache, err := helpers.InitJWKCache(ctx, &cfg)
	if err != nil {
		return nil, err
	}

	/* Initialize the Middlewares */
	dbTxMiddleware := middleware.DbTxMiddleware{
		DBClient: dbClient,
	}

	r.Use(dbTxMiddleware.HandleTransaction)

	jwtAuthMiddleware := middleware.JWTAuthMiddleware{
		Config:   &cfg,
		JwkCache: jwkCache,
	}
	orgMiddleware := middleware.OrganizationMiddleware{
		Config:                 &cfg,
		OrganizationRepository: organizationsRepository,
		RegistryRepository:     registryRepository,
	}
	extractBasicAuthMiddleware := middleware.ExtractBasicCredentialsMiddleware{
		Config:        &cfg,
		CognitoClient: helpers.GetCognitoClient(ctx, cfg),
	}

	/* Initialize HTTP Handlers */
	healthHandler := HealthHandler{}

	v2LoginHandler := V2LoginHandler{}

	v2PingHandler := V2PingHandler{}

	v2BlobsHandler := V2BlobsHandler{
		Config:              &cfg,
		S3Client:            s3Client,
		S3PresignClient:     s3Presigner,
		BlobChunkRepository: blobChunkRepository,
	}

	v2ManifestsHandlers := V2ManifestsHandler{
		Config:                &cfg,
		S3Client:              s3Client,
		S3PresignClient:       s3Presigner,
		ManifestRepository:    manifestRepository,
		RepositoryRepository:  repositoryRepository,
		ManifestTagRepository: manifestTagRepository,
	}

	v2ReferrersHandlers := V2ReferrersHandler{
		Config:               &cfg,
		S3Client:             s3Client,
		S3PresignClient:      s3Presigner,
		ManifestRepository:   manifestRepository,
		RepositoryRepository: repositoryRepository,
	}

	v2TagsHandler := V2TagsHandler{
		Config:                &cfg,
		RepositoryRepository:  repositoryRepository,
		ManifestTagRepository: manifestTagRepository,
	}

	registriesHandler := RegistriesHandler{
		Config:             &cfg,
		RegistryRepository: registryRepository,
	}

	repositoriesHandler := RepositoriesHandler{
		Config:               &cfg,
		RepositoryRepository: repositoryRepository,
	}

	artifactsHandler := ArtifactsHandler{
		Config:               &cfg,
		RepositoryRepository: repositoryRepository,
		ManifestRepository:   manifestRepository,
	}

	vulnerabilitiesHandler := VulnerabilitiesHandlers{
		Config:               &cfg,
		RepositoryRepository: repositoryRepository,
		ManifestRepository:   manifestRepository,
	}

	pythonHandler := PythonHandler{
		Config:               &cfg,
		ManifestRepository:   manifestRepository,
		RepositoryRepository: repositoryRepository,
	}

	organizationsHandlers := OrganizationsHandler{
		Config:                       &cfg,
		OrganizationRepository:       organizationsRepository,
		OrganizationInviteRepository: organizationInviteRepository,
		UserRepository:               userRepository,
		RepositoryRepository:         repositoryRepository,
		ManifestRepository:           manifestRepository,
		RegistryRepository:           registryRepository,
		SESClient:                    sesClient,
	}
	organizationInviteHandler := OrganizationInvitesHandler{
		Config:                       &cfg,
		OrganizationInviteRepository: organizationInviteRepository,
	}

	/* Configure API routes */

	// Specificy the health endpoint
	r.Get("/api/v1/health", healthHandler.GetHealth)

	// Route group for uploading and downloading python packages from a specific org and registryt
	r.Route("/api/v1/{organizationSlug:[a-z0-9-]+}/{registrySlug:[a-z0-9-]+}/python", func(pythonApiV1 chi.Router) {
		pythonApiV1.Use(extractBasicAuthMiddleware.Validate)
		pythonApiV1.Use(jwtAuthMiddleware.Validate)
		pythonApiV1.Use(orgMiddleware.ValidateOrgAndRegistry)
		pythonApiV1.Post("/", pythonHandler.UploadPythonPackage)
		pythonApiV1.Get("/simple/{packageName}", pythonHandler.SimpleRepositoryIndex)
		pythonApiV1.Get("/simple/{packageName}/{fileName}", pythonHandler.DownloadPackage)
	})

	// Route group for the administrative (non-OCI) APIs
	r.Route("/api/v1", func(authenticatedApiV1 chi.Router) {
		authenticatedApiV1.Use(jwtAuthMiddleware.Validate)
		// Route for getting invites for the authenticated user
		authenticatedApiV1.Get("/invites", organizationInviteHandler.GetInvitesForUser)

		// Route group to accept or reject a specific invite
		authenticatedApiV1.Route("/invites/{inviteId:[a-zA-Z0-9]+}", func(inviteScopedRoutes chi.Router) {
			inviteScopedRoutes.Post("/reject", organizationInviteHandler.RejectInvite)
			inviteScopedRoutes.Post("/accept", organizationInviteHandler.AcceptInvite)
		})

		// Route for getting all the organizations a user belong to
		authenticatedApiV1.Get("/organizations", organizationsHandlers.GetOrganizationsForUser)
		// Route for creating a new org
		authenticatedApiV1.Post("/organizations", organizationsHandlers.CreateOrganization)

		// Route group to interact with an organization's resources
		authenticatedApiV1.Route("/organizations/{organizationSlug:[a-z0-9-]+}", func(orgScopedRoutes chi.Router) {
			orgScopedRoutes.Use(orgMiddleware.ValidateOrg)
			// Route for getting the specified organization by slug for a user
			orgScopedRoutes.Get("/", organizationsHandlers.GetOrganizationsBySlugForUser)
			// Route for getting the organization stats
			orgScopedRoutes.Get("/stats", organizationsHandlers.GetOrganizationStats)
			// Route for getting the organization members
			orgScopedRoutes.Get("/members", organizationsHandlers.GetOrganizationMembers)
			// Route for inviting a user to the organization
			orgScopedRoutes.Post("/members", organizationsHandlers.InviteToOrganization)
			// Route for getting all organization's registries
			orgScopedRoutes.Get("/registries", registriesHandler.GetRegistries)
			// Route for creating a registry in the organization
			orgScopedRoutes.Post("/registries", registriesHandler.CreateRegistries)
			// Route for getting all artifacts in the organization
			orgScopedRoutes.Get("/artifacts", artifactsHandler.GetArtifactsForOrg)
		})

		// Route group to interact with a specific registry in the organization
		authenticatedApiV1.Route("/organizations/{organizationSlug:[a-z0-9-]+}/registries/{registrySlug:[a-z0-9-]+}", func(registryScopedRoutes chi.Router) {
			registryScopedRoutes.Use(orgMiddleware.ValidateOrgAndRegistry)
			// Route to get a registry by its slug
			registryScopedRoutes.Get("/", registriesHandler.GetRegistry)
		})

		// Route group to interact with a registry's repositories
		authenticatedApiV1.Route("/organizations/{organizationSlug:[a-z0-9-]+}/registries/{registrySlug:[a-z0-9-]+}/repositories", func(repositoryScopedRoutes chi.Router) {
			repositoryScopedRoutes.Use(orgMiddleware.ValidateOrgAndRegistry)
			// Route to get all registry repositories
			repositoryScopedRoutes.Get("/", repositoriesHandler.GetRepositories)

			// Route to get all vulnerabilities for a specific artifact by its digest
			apiCustomMux.Get(getRepositoryRegexRoute()+`\/artifacts\/(?P<digest>[\/\w:]+)/vulnerabilities`, vulnerabilitiesHandler.GetVulnerabilitiesForArtifact)
			// Route to get an artifact's details by its digest
			apiCustomMux.Get(getRepositoryRegexRoute()+`\/artifacts\/(?P<digest>[\/\w:]+)`, artifactsHandler.GetArtifactByDigest)
			// Route to get all artifacts in the repository
			apiCustomMux.Get(getRepositoryRegexRoute()+`\/artifacts`, artifactsHandler.GetArtifactsForRepository)
			// Route to get a specific repository by its slug
			apiCustomMux.Get(getRepositoryRegexRoute(), repositoriesHandler.GetRepository)
			// This exposes all the apiCustomMux routes in the repositoryScopedRoutes router
			repositoryScopedRoutes.HandleFunc("/*", apiCustomMux.Handle)
		})

	})

	// Route to login for OCI client
	r.With(extractBasicAuthMiddleware.Validate).Get("/v2/login", v2LoginHandler.Login)
	// Route group for all OCI enpoints
	r.Route("/v2", func(authenticatedOciV2 chi.Router) {
		authenticatedOciV2.Use(jwtAuthMiddleware.Validate)
		// Route for the ping endpoint
		authenticatedOciV2.Get("/", v2PingHandler.Ping)

		// Route group for interaction with OCI resources within an organization's registry
		authenticatedOciV2.Route("/{organizationSlug:[a-z-]+}/{registrySlug:[a-z-]+}", func(registryScopedOCIRoutes chi.Router) {
			registryScopedOCIRoutes.Use(orgMiddleware.ValidateOrgAndRegistry)

			// Route to initiate an upload session for a blob
			v2CustomMux.Post(getRepositoryRegexRoute()+`\/blobs\/uploads`, v2BlobsHandler.InitiateUploadSession)
			// Route to upload a blob
			v2CustomMux.Patch(getRepositoryRegexRoute()+`\/blobs\/uploads\/(?P<uploadId>[\w]+)`, v2BlobsHandler.UploadBlob)
			// Route to finalize the upload session of a blob
			v2CustomMux.Put(getRepositoryRegexRoute()+`\/blobs\/uploads\/(?P<uploadId>[\w]+)`, v2BlobsHandler.FinalizeBlobUploadSession)
			// Route to retrieve an active upload session for a blob
			v2CustomMux.Get(getRepositoryRegexRoute()+`\/blobs\/uploads\/(?P<uploadId>[\w]+)`, v2BlobsHandler.GetBlobUploadSession)
			// Route to download a blob by its digest
			v2CustomMux.Get(getRepositoryRegexRoute()+`\/blobs\/(?P<digest>[\/\w:]+)`, v2BlobsHandler.DownloadBlob)
			// Route to assert the existance of a blob by its digest
			v2CustomMux.Head(getRepositoryRegexRoute()+`\/blobs\/(?P<digest>[\/\w:]+)`, v2BlobsHandler.HeadBlob)
			// Route to delete a blob. Not supported but needed for OCI compliance
			v2CustomMux.Delete(getRepositoryRegexRoute()+`\/blobs\/(?P<digest>[\/\w:]+)`, v2BlobsHandler.DeleteBlob)

			// Route to Upload a manifest with its reference
			v2CustomMux.Put(getRepositoryRegexRoute()+`\/manifests\/(?P<reference>[\w:._-]+)`, v2ManifestsHandlers.UploadManifest)
			// Route to download a manifest by reference
			v2CustomMux.Get(getRepositoryRegexRoute()+`\/manifests\/(?P<reference>[\w:._-]+)`, v2ManifestsHandlers.DownloadManifest)
			// Route to assert the exitance of a manifest by reference
			v2CustomMux.Head(getRepositoryRegexRoute()+`\/manifests\/(?P<reference>[\w:._-]+)`, v2ManifestsHandlers.HeadManifest)
			// Route to delete a manifest
			v2CustomMux.Delete(getRepositoryRegexRoute()+`\/manifests\/(?P<reference>[\w:._-]+)`, v2ManifestsHandlers.DeleteManifestOrTag)
			// Route to get a manifest's referrers
			v2CustomMux.Get(getRepositoryRegexRoute()+`\/referrers\/(?P<reference>[\w:._-]+)`, v2ReferrersHandlers.GetReferrersForDigest)
			// Route to get a manifest's tags
			v2CustomMux.Get(getRepositoryRegexRoute()+`\/tags\/list`, v2TagsHandler.ListTags)
			// This exposes all the v2CustomMux routes in the registryScopedOCIRoutes router
			registryScopedOCIRoutes.HandleFunc("/*", v2CustomMux.Handle)
		})
	})

	return &Router{
		engine: r,
	}, nil
}

// Run the server using the specified port
func (r Router) Run(ctx context.Context, portBinding string) error {
	srv := &http.Server{
		Addr:              portBinding,
		Handler:           r.engine,
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: 3 * time.Second,
	}

	fmt.Println("Starting Server...")
	go r.listen(srv)
	<-ctx.Done()
	fmt.Println("shutting down server")

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

// Regex route to be used to to support repository names that contain slashes
func getRepositoryRegexRoute() string {
	return fmt.Sprintf(`(?P<repositoryName>%s)`, `[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*(\/[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*)*`)
}
