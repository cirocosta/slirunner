package runnable_test

import (
	"bytes"
	"context"

	"github.com/cirocosta/slirunner/runnable"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ShellCommand", func() {

	var (
		err    error
		cmd    string
		buffer bytes.Buffer

		ctx context.Context = context.TODO()
	)

	JustBeforeEach(func() {
		err = runnable.NewShellCommand(cmd, &buffer).Run(ctx)
	})

	Context("with a command that fails", func() {
		BeforeEach(func() {
			cmd = "exit 1"
		})

		It("errors", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Context("with a command that succeeds", func() {
		BeforeEach(func() {
			cmd = "exit 0"
		})

		It("doesn't error", func() {
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
