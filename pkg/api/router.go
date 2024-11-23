package api

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/helpers"
	"github.com/estebangarcia/cm3070-final-project/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	engine *chi.Mux
}

func NewRouter(ctx context.Context, cfg config.AppConfig) (*Router, error) {
	r := chi.NewRouter()

	healthHandler := HealthHandler{}
	v2LoginHandler := V2LoginHandler{
		Config:        &cfg,
		CognitoClient: helpers.GetCognitoClient(ctx, cfg),
	}
	v2PingHandler := V2PingHandler{}

	v2BlobsHandler := V2BlobsHandler{
		Config:   &cfg,
		S3Client: helpers.GetS3Client(ctx, cfg),
	}

	jwkCache, err := helpers.InitJWKCache(ctx, &cfg)
	if err != nil {
		return nil, err
	}

	jwtAuthMiddleware := middleware.JWTAuthMiddleware{
		Config:   &cfg,
		JwkCache: jwkCache,
	}

	extractBasicAuthMiddleware := middleware.ExtractBasicCredentialsMiddleware{}

	r.Get("/health", healthHandler.GetHealth)
	r.With(extractBasicAuthMiddleware.Validate).Get("/v2/login", v2LoginHandler.Login)

	r.Route("/v2", func(authenticatedOciV2 chi.Router) {
		authenticatedOciV2.Use(jwtAuthMiddleware.Validate)
		authenticatedOciV2.Get("/", v2PingHandler.Ping)
		authenticatedOciV2.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uri := chi.URLParam(r, "*")

			routes := map[string][]struct {
				Pattern *regexp.Regexp
				Handler http.HandlerFunc
			}{
				http.MethodPost: {
					{
						Pattern: regexp.MustCompile(`^(?P<imageName>[\/\w]+)\/blobs\/uploads`), Handler: v2BlobsHandler.InitiateUploadSession,
					},
				},
				http.MethodHead: {
					{
						Pattern: regexp.MustCompile(`^(?P<imageName>[\/\w]+)\/blobs\/(?P<digest>[\/\w:]+)`), Handler: v2BlobsHandler.HeadBlob,
					},
				},
			}

			methodRoutes, exists := routes[r.Method]
			if !exists {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			for _, route := range methodRoutes {
				matches := route.Pattern.FindStringSubmatch(uri)
				if len(matches) > 0 {
					ctx := r.Context()
					for i, name := range route.Pattern.SubexpNames() {
						if i != 0 && name != "" {
							ctx = context.WithValue(ctx, name, matches[i])
						}
					}
					route.Handler(w, r.WithContext(ctx))
					return
				}
			}

			http.Error(w, "Not Found", http.StatusNotFound)
		}))
		//authenticatedOciV2.Post("/{imageName}/*/blobs/uploads/", v2BlobUploadHandler.InitiateUploadSession)

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
