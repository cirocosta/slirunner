package runnable

import (
	"context"

	"github.com/pkg/errors"
)

// Sequentially serially runs a set of runnables, cancelling the
// execution of all of the next ones once a failure occurs.
//
type Sequentially struct {
	runnables []Runnable
}

var _ Runnable = &Sequentially{}

func NewSequentially(runnables []Runnable) *Sequentially {
	return &Sequentially{
		runnables: runnables,
	}
}

func (r *Sequentially) Run(ctx context.Context) (err error) {
	for _, runnable := range r.runnables {
		err = runnable.Run(ctx)
		if err != nil {
			err = errors.Wrapf(err,
				"failed while sequentially running group of runnables")
			return
		}
	}

	return
}
