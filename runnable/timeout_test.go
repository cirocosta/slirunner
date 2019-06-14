package runnable_test

import (
	"context"
	"time"

	"github.com/cirocosta/slirunner/runnable"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// sleepingRunnable is a dummy runnable that sleeps forever until
// cancelled by its context.
//
type sleepingRunnable struct{}

func (r *sleepingRunnable) Run(ctx context.Context) (err error) {
	ch := make(chan struct{})

	go func() {
		ch <- (<-ctx.Done())
	}()

	<-ch
	return
}

var _ = Describe("WithTimeout", func() {

	JustBeforeEach(func() {
		runnable.NewWithTimeout(
			&sleepingRunnable{}, 200*time.Millisecond,
		).Run(context.TODO())
	})

	Context("with a duration specified", func() {

		Context("a runnable that would sleep forever", func() {

			It("doesnt sleep forever", func() {
				Expect(true).To(BeTrue())
			})

		})

	})

})
