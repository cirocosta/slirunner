package runnable

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Concurrently concurrently runs a set of runnables, cancelling them all at
// once in case of context cancellation.
//
type Concurrently struct {
	runnables []Runnable
}

func NewConcurrently(runnables []Runnable) *Concurrently {
	return &Concurrently{
		runnables: runnables,
	}
}

func (r *Concurrently) Run(ctx context.Context) (err error) {
	g, ctx := errgroup.WithContext(ctx)

	for _, runnable := range r.runnables {
		runnable := runnable // closure

		g.Go(func() error {
			return runnable.Run(ctx)
		})
	}

	err = g.Wait()
	if err != nil {
		err = errors.Wrapf(err,
			"failed while concurrently running group of runnables")
		return
	}

	return
}
