package runnable_test

import (
	"bytes"
	"context"
	"time"

	"github.com/cirocosta/slirunner/runnable"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ShellCommand", func() {

	var (
		err    error
		cmd    string
		buffer bytes.Buffer
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.TODO()
		buffer.Reset()
	})

	JustBeforeEach(func() {
		err = runnable.NewShellCommand(cmd, &buffer).Run(ctx)
	})

	Context("with a command that fails", func() {
		BeforeEach(func() {
			cmd = `echo "damn" && exit 1`
		})

		It("errors", func() {
			Expect(err).To(HaveOccurred())
		})

		It("writes output to `stderr`", func() {
			Expect(buffer.Bytes()).ToNot(BeEmpty())
			Expect(string(buffer.Bytes())).To(ContainSubstring("damn"))
		})
	})

	Context("with a command that succeeds", func() {
		BeforeEach(func() {
			cmd = "exit 0"
		})

		It("doesn't error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("doesn't write to `stderr`", func() {
			Expect(buffer.Bytes()).To(BeEmpty())
		})
	})

	Context("with a command that hangs forever", func() {
		BeforeEach(func() {
			cmd = "sleep 33d"
			ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
		})

		Context("cancelling it", func() {
			It("fails", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
