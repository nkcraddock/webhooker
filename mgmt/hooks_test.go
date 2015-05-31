package mgmt_test

import (
	"net/http"

	"gopkg.in/redis.v3"

	"github.com/gorilla/mux"
	"github.com/nkcraddock/webhooker/db"
	"github.com/nkcraddock/webhooker/webhooks"
	"github.com/nkcraddock/webhooker/mgmt"
	"github.com/nkcraddock/webhooker/testhelp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("hooks handler tests", func() {
	var s *testhelp.TestServer
	var store webhooks.Store
	var client *redis.Client
	regex_guid := `[\da-zA-Z]{8}-([\da-zA-Z]{4}-){3}[\da-zA-Z]{12}`
	testdata, _ := testhelp.LoadTestData("testdata/integration-tests.json")

	BeforeEach(func() {

		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			PoolSize: 10,
			DB:       10,
		})

		client.FlushDb()

		store = db.RedisHookerStore(func() *redis.Client {
			return client
		})

		h := mgmt.NewHooksHandler(store)
		m := mux.NewRouter()
		h.RegisterRoutes(m)

		s = &testhelp.TestServer{m}
	})

	AfterEach(func() {
		client.Close()
	})

	Context("POST", func() {
		It("returns created status", func() {
			res := s.POST("/hooks", testdata["h1"], nil)
			Ω(res.Code).Should(Equal(http.StatusCreated))
		})

		It("returns the location header", func() {
			res := s.POST("/hooks", testdata["h1"], nil)
			location := res.Header().Get("Location")
			Ω(location).Should(MatchRegexp(regex_guid))
		})

		It("returns the new hook", func() {
			hook := new(webhooks.Hook)
			s.POST("/hooks", testdata["h1"], hook)
			Ω(hook.Id).Should(MatchRegexp(regex_guid))
		})
	})

	Context("GET", func() {
		It("gets a list of all hooks", func() {
			s.POST("/hooks", testdata["h1"], nil)
			s.POST("/hooks", testdata["h2"], nil)

			var hooks []*webhooks.Hook
			res := s.GET("/hooks", &hooks)
			Ω(res.Code).Should(Equal(http.StatusOK))
			Ω(hooks).Should(HaveLen(2))
		})

		It("gets an individual hook", func() {
			savedhook := new(webhooks.Hook)
			s.POST("/hooks", testdata["h1"], savedhook)

			hook := new(webhooks.Hook)
			res := s.GET("/hooks/"+savedhook.Id, hook)
			Ω(res.Code).Should(Equal(http.StatusOK))
			Ω(hook.RatePerMinute).Should(Equal(savedhook.RatePerMinute))
		})
	})
})
