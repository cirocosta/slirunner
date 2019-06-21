package runnable_test

import (
	"context"
	"time"

	"github.com/cirocosta/slirunner/runnable"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
