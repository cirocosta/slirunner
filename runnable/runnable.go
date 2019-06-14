package runnable

import (
	"context"
)

// Runnable represents something that has to be executed and
// cancellable through context cancellation.
//
type Runnable interface {
	// Run runs what needed to be ran until cancelled by
	// `ctx` or the completion of its run.
	//
	// ps.: runs terminated by context cancellation or
	//      deadlines are considered not succeeded.
	//
	Run(ctx context.Context) (succeeded bool)
}
