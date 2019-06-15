package probes_test

import (
	"github.com/cirocosta/slirunner/probes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("All", func() {

	Describe("FormatProbe", func() {

		Context("having the configuration passed", func() {

			var (
				str    string
				config probes.Config
				res    string
			)

			BeforeEach(func() {
				str = "pn={{ .Pipeline }} tgt={{ .Target }}"
				config = probes.Config{
					Pipeline: "pn",
					Target:   "tgt",
				}
			})

			JustBeforeEach(func() {
				res = probes.FormatProbe(str, config)
			})

			It("formats", func() {
				Expect(res).To(Equal("pn=pn tgt=tgt"))
			})

		})

	})

})
