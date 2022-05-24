package platform

import (
	"log"
	"securebanking-test-data-initializer/pkg/common"
	"strconv"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestConsumer(t *testing.T) {
	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		Consumer: "initializer",
		Provider: "am",
		Host:     "localhost",
	}
	defer pact.Teardown()

	// Set up our expected interactions.
	pact.
		AddInteraction().
		Given("Server info exists").
		UponReceiving("A request to get Server info").
		WithRequest(dsl.Request{
			Method:  "GET",
			Path:    dsl.String("/am/json/serverinfo/*"),
			Headers: dsl.MapMatcher{"Accept": dsl.String("application/json")},
		}).
		WillRespondWith(dsl.Response{
			Status:  200,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body:    dsl.MapMatcher{"cookieName": dsl.String("iPlanetDirectory")},
		})

	var test = func() (err error) {
		common.Config.Hosts.IdentityPlatformFQDN = "localhost:" + strconv.Itoa(pact.Server.Port)
		common.Config.Hosts.Scheme = "http"
		cookie := GetCookieNameFromAm()

		assert.Equal(t, "iPlanetDirectory", cookie)
		return nil
	}
	// Run the test, verify it did what we expected and capture the contract
	if err := pact.Verify(test); err != nil {
		log.Fatalf("Error on Verify: %v", err)
	}
}
