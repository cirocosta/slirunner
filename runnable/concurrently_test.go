package runnable_test

import (
	"context"
	"errors"
	"time"

	"github.com/cirocosta/slirunner/runnable"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Concurrently", func() {

	var (
		runnables []runnable.Runnable
		cancel    context.CancelFunc
		err       error
	)

	AfterEach(func() {
		cancel()
	})

	Context("having several runnables that finish", func() {

		JustBeforeEach(func() {
			var ctx context.Context

			ctx, cancel = context.WithCancel(context.Background())
			err = runnable.NewConcurrently(runnables).Run(ctx)
		})

		BeforeEach(func() {
			runnables = []runnable.Runnable{
				&dummyRunnable{},
				&dummyRunnable{},
				&dummyRunnable{},
			}
		})

		// ¯\_(ツ)_/¯
		It("eventually runs them all", func() {
			Expect(err).NotTo(HaveOccurred())

			for _, runnable := range runnables {
				dr := runnable.(*dummyRunnable)
				Expect(dr.ran).To(BeTrue())
			}
		})

		Context("having at least one runnable that fail", func() {

			BeforeEach(func() {
				runnables = []runnable.Runnable{
					&dummyRunnable{},
					&dummyRunnable{err: errors.New("lol")},
					&dummyRunnable{},
				}
			})

			It("propagates the failure", func() {
				Expect(err).To(HaveOccurred())
			})

		})

	})

	Context("having a runnables that wait for ctx cancellation", func() {

		JustBeforeEach(func() {
			var ctx context.Context

			ctx, cancel = context.WithCancel(context.Background())
			go func() {
				err = runnable.NewConcurrently(runnables).Run(ctx)
			}()
		})

		BeforeEach(func() {
			runnables = []runnable.Runnable{
				&sleepingRunnable{},
			}
		})

		Context("and one that fails", func() {

			BeforeEach(func() {
				runnables = append(runnables, &dummyRunnable{
					err: errors.New("lol"),
				})
			})

			It("doesn't cancel them all", func() {
				Consistently(func() bool {
					r := runnables[0].(*sleepingRunnable)
					return r.finished
				}, 1*time.Second).Should(BeFalse())
			})

		})
	})

})
