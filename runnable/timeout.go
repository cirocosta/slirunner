package runnable

import (
	"context"
	"time"
)

// WithTimeout wraps a Runnable with a timeout, ensuring that it doesn't
// run forever.
//
type WithTimeout struct {
	runnable Runnable
	duration time.Duration
}

var _ Runnable = &WithTimeout{}

// NewWithTimeout instantiates a runnable that has the context automatically
// cancelled after a given duration.
//
func NewWithTimeout(runnable Runnable, duration time.Duration) *WithTimeout {
	return &WithTimeout{
		runnable: runnable,
		duration: duration,
	}
}

// Run runs the registered runnable for a maximum of `duration`, cancelling
// the context thereafter.
//
func (r *WithTimeout) Run(ctx context.Context) (err error) {
	ctx, _ = context.WithTimeout(ctx, r.duration)
	err = r.runnable.Run(ctx)
	return
}
