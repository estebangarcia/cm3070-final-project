package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"golang.org/x/sync/errgroup"
)

type PeriodicWorkerDispatcher struct {
	RunEvery time.Duration
	Group    *errgroup.Group
	dBClient *ent.Client
}

func NewPeriodicWorkerDispatcher(runEvery time.Duration) *PeriodicWorkerDispatcher {
	return &PeriodicWorkerDispatcher{
		RunEvery: runEvery,
		Group:    &errgroup.Group{},
	}
}

func (w *PeriodicWorkerDispatcher) Start(ctx context.Context, worker PeriodicWorker) {
	w.Group.Go(func() error {
		for {
			tx, err := w.dBClient.Tx(ctx)
			if err != nil {
				return err
			}

			if err := worker.Handle(ctx); err != nil {
				fmt.Println(err)
				tx.Rollback()
				return err
			}
			if err := tx.Commit(); err != nil {
				fmt.Printf("error commiting transaction %v", err)
				return err
			}
			select {
			case <-time.After(w.RunEvery):
				continue
			case <-ctx.Done():
				log.Printf("exiting %T", worker)
				return nil
			}
		}
	})
}
