package runnable_test

import (
	"context"
	"errors"

	"github.com/cirocosta/slirunner/runnable"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Concurrently", func() {

	var (
		runnables []runnable.Runnable
		err       error
	)

	Describe("with propagation", func() {

		JustBeforeEach(func() {
			err = runnable.NewConcurrently(runnables).Run(context.TODO())
		})

		Context("having several runnables that finish", func() {

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

		})

		Context("having some runnables that sleep forever", func() {

			BeforeEach(func() {
				runnables = []runnable.Runnable{
					&sleepingRunnable{},
					&sleepingRunnable{},
				}
			})

			Context("but one that fails", func() {

				BeforeEach(func() {
					runnables = append(runnables, &dummyRunnable{
						err: errors.New("lol"),
					})
				})

				It("propagates the failure", func() {
					Expect(err).To(HaveOccurred())
				})

				It("cancels them all", func() {
					// otherwise, we'd not even run this assertion
					Expect(true).To(BeTrue())
				})

			})
		})

	})

	Describe("WithoutErrorPropagation", func() {

		JustBeforeEach(func() {
			err = runnable.NewConcurrentlyWithoutErrorPropagation(runnables).Run(context.TODO())
		})

		Context("having at least one runnable that fail", func() {

			BeforeEach(func() {
				runnables = []runnable.Runnable{
					&dummyRunnable{},
					&dummyRunnable{err: errors.New("lol")},
					&dummyRunnable{},
				}
			})

			It("doesnt propagate the failure", func() {
				Expect(err).NotTo(HaveOccurred())
			})

		})

	})

})
