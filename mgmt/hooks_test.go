package mgmt_test

import (
	"net/http"

	"github.com/nkcraddock/mervis/testhelp"
	"github.com/nkcraddock/webhooker/mgmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("hooks handler tests", func() {
	var s *testhelp.TestServer
	BeforeEach(func() {
		handlers := []mgmt.Handler{mgmt.NewHooksHandler()}
		m, _ := mgmt.NewMgmtServer(handlers)

		s = &testhelp.TestServer{
			Handler: m,
		}
	})
	Context("Hooks resource", func() {
		It("", func() {
			res := s.GET("/api/hooks")
			Ω(res.Code).Should(Equal(http.StatusOK))
		})
	})
})
