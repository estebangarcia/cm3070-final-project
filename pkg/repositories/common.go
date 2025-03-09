package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
)

func getClient(ctx context.Context) *ent.Client {
	return ctx.Value("dbClient").(*ent.Client)
}
