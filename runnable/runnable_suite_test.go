package runnable_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// sleepingRunnable is a dummy runnable that sleeps forever until
// cancelled by its context.
//
type sleepingRunnable struct {
	finished bool
}

func (r *sleepingRunnable) Run(ctx context.Context) (err error) {
	ch := make(chan struct{})

	go func() {
		ch <- (<-ctx.Done())
	}()

	<-ch
	r.finished = true
	return
}

// dummyRunnable implements the Runnable interface, capturing how
// many times it ran.
//
type dummyRunnable struct {
	ran bool
	err error
}

func (r *dummyRunnable) Run(ctx context.Context) (err error) {
	r.ran = true
	err = r.err
	return
}

func TestRunnable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Runnable Suite")
}
