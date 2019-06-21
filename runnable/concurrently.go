package runnable

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Concurrently concurrently runs a set of runnables, cancelling them all at
// once in case of context cancellation.
//
type Concurrently struct {
	runnables []Runnable
}

var _ Runnable = &Concurrently{}

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

// ConcurrentlyWithoutErrorPropagation runs a set of runnables without
// ever caring about the errors that those runnables might return.
//
// It just runs runnables and forget about what happened.
//
type ConcurrentlyWithoutErrorPropagation struct {
	runnables []Runnable
}

func NewConcurrentlyWithoutErrorPropagation(runnables []Runnable) *ConcurrentlyWithoutErrorPropagation {
	return &ConcurrentlyWithoutErrorPropagation{
		runnables: runnables,
	}
}

func (r *ConcurrentlyWithoutErrorPropagation) Run(ctx context.Context) (err error) {
	var wg sync.WaitGroup

	for _, runnable := range r.runnables {
		wg.Add(1)
		go func(r Runnable) {
			runnable.Run(ctx)
			wg.Done()
		}(runnable)
	}

	wg.Wait()
	return
}
