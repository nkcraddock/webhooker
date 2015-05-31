package mgmt_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mgmt server tests", func() {
	Context("afsd", func() {
		It("hehes", func() {
			Ω(1).ShouldNot(Equal(2))
		})
	})
})
