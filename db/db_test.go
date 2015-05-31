package db_test

import (
	"testing"

	"gopkg.in/redis.v3"

	"github.com/nkcraddock/webhooker/db"
	"github.com/nkcraddock/webhooker/domain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "db tests")
}

var _ = Describe("RedisHookerStore integration tests", func() {
	var store domain.Store
	var client *redis.Client
	var testhook *domain.Hook

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

		testhook = domain.NewHook("url", 100)
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
			one := domain.NewHook("one", 1)
			two := domain.NewHook("two", 2)
			store.SaveHook(one)
			store.SaveHook(two)

			hooks, err := store.GetHooks("")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(hooks).Should(HaveLen(2))
		})
	})
})
