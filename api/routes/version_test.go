package routes_test

import (
	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/routes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("getRequestResult", func() {

	expected := map[string]string{
		"version": apiContext.GetAPIVersion(),
		"date":    apiContext.GetAPIReleaseDate(),
	}

	apiContext.SetOnceConfig()
	config := apiContext.APIConfiguration

	Context("When version and date are requested", func() {
		It("Should return a map with API version and date", func() {
			Expect(routes.GetRequestResult(config)).To(Equal(expected))
		})
	})

})
