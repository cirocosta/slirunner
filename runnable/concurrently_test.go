package runnable_test

import (
	"context"

	"github.com/cirocosta/slirunner/runnable"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// dummyRunnable implements the Runnable interface, capturing how
// many times it ran.
//
type dummyRunnable struct {
	ran bool
}

func (r *dummyRunnable) Run(ctx context.Context) (err error) {
	r.ran = true
	return
}

var _ = Describe("Concurrently", func() {

	var runnables []runnable.Runnable

	JustBeforeEach(func() {
		err := runnable.NewConcurrently(runnables).Run(context.TODO())
		Expect(err).ToNot(HaveOccurred())
	})

	Context("having several runnables", func() {

		BeforeEach(func() {
			runnables = []runnable.Runnable{
				&dummyRunnable{},
				&dummyRunnable{},
				&dummyRunnable{},
			}
		})

		// ¯\_(ツ)_/¯
		It("eventually runs them all", func() {
			for _, runnable := range runnables {
				dr := runnable.(*dummyRunnable)
				Expect(dr.ran).To(BeTrue())
			}
		})
	})

})
