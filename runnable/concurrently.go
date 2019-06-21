package runnable

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Concurrently runs a set of runnables without
// ever caring about the errors that those runnables might return.
//
// It just runs runnables and forget about what happened.
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
	err = runConcurrently(ctx, r.runnables, new(errgroup.Group))
	return
}

// ConcurrentlyFailFast concurrently runs a set of runnables, cancelling them all at
// once in case of context cancellation.
//
type ConcurrentlyFailFast struct {
	runnables []Runnable
}

var _ Runnable = &ConcurrentlyFailFast{}

func NewConcurrentlyFailFast(runnables []Runnable) *ConcurrentlyFailFast {
	return &ConcurrentlyFailFast{
		runnables: runnables,
	}
}

func (r *ConcurrentlyFailFast) Run(ctx context.Context) (err error) {
	var g *errgroup.Group

	g, ctx = errgroup.WithContext(ctx)
	err = runConcurrently(ctx, r.runnables, g)
	return
}

// runConcurrently runs a bunch of runnables concurrently.
//
func runConcurrently(ctx context.Context, runnables []Runnable, g *errgroup.Group) (err error) {
	for _, runnable := range runnables {
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
