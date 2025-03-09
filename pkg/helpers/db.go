package helpers

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"

	_ "github.com/lib/pq"
)

func GetDBClient(ctx context.Context, cfg *config.AppConfig) (*ent.Client, error) {
	drv, err := entsql.Open("postgres", cfg.Database.DSN)
	if err != nil {
		return nil, err
	}

	client := ent.NewClient(ent.Driver(drv))

	if cfg.Database.Debug {
		sqlDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
			for _, inter := range i {
				fmt.Printf("%v\n", inter)
			}
		})

		client = ent.NewClient(ent.Driver(sqlDrv)).Debug()
	}

	return client, nil
}
