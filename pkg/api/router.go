package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/helpers"
	"github.com/estebangarcia/cm3070-final-project/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter(ctx context.Context, cfg config.AppConfig) (*Router, error) {
	r := gin.Default()

	healthHandler := HealthHandler{}
	v2LoginHandler := V2LoginHandler{}
	v2PingHandler := V2PingHandler{}

	jwkCache, err := helpers.InitJWKCache(ctx, &cfg)
	if err != nil {
		return nil, err
	}

	authMiddleware := middleware.AuthMiddleware{
		Config:   &cfg,
		JwkCache: jwkCache,
	}

	r.GET("/health", healthHandler.GetHealth)

	r.GET("/v2/login", v2LoginHandler.Login)

	authenticatedOciV2 := r.Group("/v2", authMiddleware.ValidateJWT)
	{
		authenticatedOciV2.GET("/", v2PingHandler.Ping)
	}

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
