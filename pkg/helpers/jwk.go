package helpers

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

// Initialize the JWK cache, this is used to cache the JWK public keys
// used to verify the signature of a JWT token
func InitJWKCache(ctx context.Context, cfg *config.AppConfig) (*jwk.Cache, error) {
	cache, err := jwk.NewCache(ctx, httprc.NewClient())
	if err != nil {
		return nil, err
	}

	if err := cache.Register(ctx, cfg.GetCognitoJWKUrl()); err != nil {
		return nil, err
	}

	return cache, nil
}
