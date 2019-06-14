package runnable_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRunnable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Runnable Suite")
}
