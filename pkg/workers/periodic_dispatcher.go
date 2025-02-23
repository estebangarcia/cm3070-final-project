package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

type PeriodicWorkerDispatcher struct {
	RunEvery time.Duration
	Group    *errgroup.Group
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
			err := worker.Handle(ctx)
			if err != nil {
				fmt.Println(err)
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
