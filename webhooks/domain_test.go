package webhooks_test

import (
	"testing"

	"github.com/nkcraddock/webhooker/webhooks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "webhooks tests")
}

var _ = Describe("Domain unit tests", func() {
	Context("NewHook", func() {
		It("should initialize an Id and Secret", func() {
			hook := webhooks.NewHook("url", 5)
			Ω(hook.Id).ShouldNot(Equal(""))
			Ω(hook.Secret).ShouldNot(Equal(""))
		})
	})

	Context("Hook.NewFilter", func() {
		It("creates a new filter for the hook", func() {
			hook := webhooks.NewHook("url", 5)
			filter := hook.NewFilter("src", "evt", "key")
			Ω(filter.Hook).Should(Equal(hook.Id))
		})
	})
})
