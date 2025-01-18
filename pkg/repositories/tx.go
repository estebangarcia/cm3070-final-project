package repositories

import (
	"context"
	"fmt"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
)

type TransactionalRepository struct {
	dbClient *ent.Client
}

func (tr *TransactionalRepository) WithTx(ctx context.Context, fn func(tx *ent.Tx) error) error {
	tx, err := tr.dbClient.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
