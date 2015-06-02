package db_test

import (
	"testing"

	"gopkg.in/redis.v3"

	"github.com/nkcraddock/webhooker/db"
	"github.com/nkcraddock/webhooker/webhooks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "db tests")
}

var _ = Describe("RedisHookerStore integration tests", func() {
	var store webhooks.Store
	var client *redis.Client
	var testhook *webhooks.Hook
	var testfilter *webhooks.Filter

	BeforeEach(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			PoolSize: 10,
			DB:       9,
		})

		client.FlushDb()

		store = db.RedisHookerStore(func() *redis.Client {
			return client
		})

		testhook = webhooks.NewHook("url", 100)
		testfilter = testhook.NewFilter("testfilter", "evt", "key")
	})

	AfterEach(func() {
		client.Close()
	})

	Context("SaveHook", func() {
		It("saves a new hook", func() {
			err := store.SaveHook(testhook)
			Ω(err).ShouldNot(HaveOccurred())

			storedhook, err := store.GetHook(testhook.Id)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(storedhook).ShouldNot(BeNil())
			Ω(storedhook.Id).Should(Equal(testhook.Id))
		})

		It("updates an existing hook", func() {
			// add testhook
			store.SaveHook(testhook)
			// change a value
			testhook.Secret = "haha"
			// save the change
			store.SaveHook(testhook)

			storedhook, _ := store.GetHook(testhook.Id)
			Ω(storedhook.Id).Should(Equal(testhook.Id))
			Ω(storedhook.Secret).Should(Equal("haha"))
		})
	})

	Context("GetHook", func() {
		It("retrieves a hook from redis", func() {
			store.SaveHook(testhook)
			storedhook, err := store.GetHook(testhook.Id)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(storedhook).ShouldNot(BeNil())
			Ω(storedhook).Should(Equal(testhook))
		})
	})

	Context("GetHooks", func() {
		It("lists all hooks", func() {
			one := webhooks.NewHook("one", 1)
			two := webhooks.NewHook("two", 2)
			store.SaveHook(one)
			store.SaveHook(two)

			hooks, err := store.GetHooks("")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(hooks).Should(HaveLen(2))
		})
	})

	Context("SaveFilter", func() {
		It("saves a new filter", func() {
			store.SaveHook(testhook)
			filter := testhook.NewFilter("url", "evt", "key")
			err := store.SaveFilter(filter)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("GetFilters", func() {
		It("returns the filters for a given hook id", func() {
			hookOne := webhooks.NewHook("one", 1)
			store.SaveHook(hookOne)
			filterOne := hookOne.NewFilter("fone", "evt", "key")
			store.SaveFilter(filterOne)

			hookTwo := webhooks.NewHook("two", 2)
			store.SaveHook(hookTwo)
			filterTwo := hookTwo.NewFilter("ftwo", "evt", "key")
			store.SaveFilter(filterTwo)

			filters, err := store.GetFilters(hookOne.Id)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(filters).Should(HaveLen(1))
			Ω(filters[0]).Should(Equal(filterOne))

			filters, err = store.GetFilters(hookTwo.Id)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(filters).Should(HaveLen(1))
			Ω(filters[0]).Should(Equal(filterTwo))
		})
	})

	Context("DeleteFilter", func() {
		It("deletes a filter", func() {
			store.SaveHook(testhook)
			store.SaveFilter(testfilter)

			err := store.DeleteFilter(testfilter.Hook, testfilter.Id)
			Ω(err).ShouldNot(HaveOccurred())

			filter, _ := store.GetFilters(testfilter.Hook)
			Ω(filter).Should(HaveLen(0))
		})
	})

})
