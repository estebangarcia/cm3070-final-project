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
	// Open DB connection to the configured DSN
	drv, err := entsql.Open("postgres", cfg.Database.DSN)
	if err != nil {
		return nil, err
	}

	client := ent.NewClient(ent.Driver(drv))

	// If the configuration sets the debug to true then we configure the driver to
	// output executed queries to the terminal
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
