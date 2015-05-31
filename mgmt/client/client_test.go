package client_test

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nkcraddock/webhooker/mgmt/client"
	"github.com/nkcraddock/webhooker/testhelp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Resource Tests")
}

var _ = Describe("client.handler", func() {
	Context("", func() {
		var srv *testhelp.TestServer
		var c *client.Handler

		BeforeEach(func() {
			router := mux.NewRouter()
			srv = &testhelp.TestServer{router}
			c = client.NewHandler(&client.BinDataLocator{})
			c.RegisterRoutes(router)
		})

		It("Serves up data from bindata", func() {
			res := srv.GET("/", nil)
			Ω(res.Code).Should(Equal(http.StatusOK))
		})
	})
})
