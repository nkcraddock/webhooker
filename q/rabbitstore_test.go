package q_test

import (
	"testing"

	"github.com/michaelklishin/rabbit-hole"
	"github.com/nkcraddock/webhooker/q"
	"github.com/nkcraddock/webhooker/webhooks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "q tests")
}

const (
	username = "guest"
	password = "guest"
	uri      = "http://localhost:15672"
	vh       = "test"
)

var _ = Describe("rabbitstore", func() {
	var store *q.RabbitStore
	var rmq *rabbithole.Client

	BeforeEach(func() {
		var err error
		rmq, err = rabbithole.NewClient(uri, username, password)
		Ω(err).ShouldNot(HaveOccurred())

		rmq.DeleteVhost(vh)
		store = q.NewRabbitStore(rmq, vh)
	})

	Context("SaveHook", func() {
		It("creates a queue, an exchange, and binds them", func() {
			hook := webhooks.NewHook("test", 5)
			err := store.SaveHook(hook)
			Ω(err).ShouldNot(HaveOccurred())

			_, err = rmq.GetQueue(vh, hook.Id)
			Ω(err).ShouldNot(HaveOccurred())

			_, err = rmq.GetExchange(vh, hook.Id)
			Ω(err).ShouldNot(HaveOccurred())

			bindings, err := rmq.ListQueueBindings(vh, hook.Id)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(bindings).Should(HaveLen(2))
		})
	})

	Context("GetHook", func() {
		It("gets a single hook", func() {
			hook1 := webhooks.NewHook("one", 1)
			store.SaveHook(hook1)

			hook, err := store.GetHook(hook1.Id)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(hook.Id).Should(Equal(hook.Id))
		})
	})

	Context("GetHooks", func() {
		It("gets a list of all hooks", func() {
			hook1 := webhooks.NewHook("one", 1)
			hook2 := webhooks.NewHook("two", 2)
			store.SaveHook(hook1)
			store.SaveHook(hook2)

			hooks, err := store.GetHooks("")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(hooks).Should(HaveLen(2))
		})
	})

	Context("SaveFilter", func() {
		It("adds a new filter", func() {
			hook1 := webhooks.NewHook("one", 1)
			store.SaveHook(hook1)

			f1 := hook1.NewFilter("a", "b", "c")
			err := store.SaveFilter(f1)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(f1.RmqProps).ShouldNot(Equal(""))
		})
	})

})
