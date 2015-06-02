package mgmt_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/nkcraddock/webhooker/db"
	"github.com/nkcraddock/webhooker/mgmt"
	"github.com/nkcraddock/webhooker/testhelp"
	"github.com/nkcraddock/webhooker/webhooks"
	"gopkg.in/redis.v3"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mgmt tests")
}

var _ = Describe("mgmt api integration tests", func() {
	var s *testhelp.TestServer
	var store webhooks.Store
	var client *redis.Client
	var testdata testhelp.TestData
	regex_guid := `[\da-zA-Z]{8}-([\da-zA-Z]{4}-){3}[\da-zA-Z]{12}`

	BeforeEach(func() {
		testdata, _ = testhelp.LoadTestData("testdata/integration-tests.json")
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
		f := mgmt.NewFiltersHandler(store)
		m := mux.NewRouter()
		h.RegisterRoutes(m)
		f.RegisterRoutes(m)
		s = &testhelp.TestServer{m}
	})

	AfterEach(func() {
		client.Close()
	})

	// hooks
	Context("hooks", func() {
		Context("POST /hooks", func() {
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
				Ω(hook.CallbackUrl).Should(Equal(testdata["h1"]["url"]))
			})

		})

		Context("GET /hooks", func() {
			It("gets a list of all hooks", func() {
				s.POST("/hooks", testdata["h1"], nil)
				s.POST("/hooks", testdata["h2"], nil)

				var hooks []*webhooks.Hook
				res := s.GET("/hooks", &hooks)
				Ω(res.Code).Should(Equal(http.StatusOK))
				Ω(hooks).Should(HaveLen(2))
			})

			It("doesnt return the secrets", func() {
				s.POST("/hooks", testdata["h1"], nil)
				s.POST("/hooks", testdata["h2"], nil)

				var hooks []map[string]interface{}
				s.GET("/hooks", &hooks)

				for _, h := range hooks {
					Ω(h).ShouldNot(HaveKey("secret"))
				}
			})
		})

		Context("GET /hooks/:id", func() {
			It("gets an individual hook", func() {
				savedhook := new(webhooks.Hook)
				s.POST("/hooks", testdata["h1"], savedhook)

				hook := new(webhooks.Hook)
				res := s.GET("/hooks/"+savedhook.Id, hook)

				Ω(res.Code).Should(Equal(http.StatusOK))
				Ω(hook.Id).Should(Equal(savedhook.Id))
				Ω(hook.RatePerMinute).Should(Equal(savedhook.RatePerMinute))
			})

			It("doesnt return the secret", func() {
				savedhook := new(webhooks.Hook)
				s.POST("/hooks", testdata["h1"], savedhook)

				var hook map[string]interface{}
				s.GET("/hooks/"+savedhook.Id, &hook)
				Ω(hook).ShouldNot(HaveKey("secret"))
			})
		})
	}) // - hooks

	// filters
	Context("filters", func() {
		newhook := new(webhooks.Hook)
		urifmt := "/hooks/%s/filters"

		BeforeEach(func() {
			// insert the test hook
			s.POST("/hooks", testdata["h1"], newhook)
		})

		Context("POST /hooks/:id/filters", func() {
			It("returns a created status", func() {
				testdata["f1"]["hook"] = newhook.Id
				uri := fmt.Sprintf(urifmt, testdata["f1"]["hook"])

				savedFilter := new(webhooks.Filter)
				res := s.POST(uri, testdata["f1"], savedFilter)

				Ω(res.Code).Should(Equal(http.StatusCreated))
			})

			It("requires a valid hook", func() {
				testdata["f1"]["hook"] = "goobledy"
				uri := fmt.Sprintf(urifmt, testdata["f1"]["hook"])

				res := s.POST(uri, testdata["f1"], nil)
				Ω(res.Code).Should(Equal(http.StatusNotFound))
			})

			It("returns the new filter", func() {
				testdata["f1"]["hook"] = newhook.Id
				uri := fmt.Sprintf(urifmt, testdata["f1"]["hook"])

				savedFilter := new(webhooks.Filter)
				s.POST(uri, testdata["f1"], savedFilter)
				Ω(savedFilter.Id).Should(MatchRegexp(regex_guid))
			})

			It("returns the location header", func() {
				testdata["f1"]["hook"] = newhook.Id
				uri := fmt.Sprintf(urifmt, testdata["f1"]["hook"])

				res := s.POST(uri, testdata["f1"], nil)
				Ω(res.Header()).Should(HaveKey("Location"))
				Ω(res.Header().Get("Location")).Should(MatchRegexp(regex_guid))
			})
		})

		Context("GET /hooks/:hook/filters", func() {
			It("returns a single filter", func() {
				testdata["f1"]["hook"] = newhook.Id
				uri := fmt.Sprintf(urifmt, testdata["f1"]["hook"])
				savedFilter := new(webhooks.Filter)
				s.POST(uri, testdata["f1"], savedFilter)

				var filters []*webhooks.Filter
				res := s.GET(uri, &filters)
				Ω(res.Code).Should(Equal(http.StatusOK))
				Ω(filters).Should(HaveLen(1))
				Ω(filters[0]).Should(Equal(savedFilter))
			})
		})
	})
})
