package helpers

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"

	_ "github.com/lib/pq"
)

func GetDBClient(ctx context.Context, cfg *config.AppConfig) (*ent.Client, error) {
	return ent.Open("postgres", cfg.Database.DSN)
}
