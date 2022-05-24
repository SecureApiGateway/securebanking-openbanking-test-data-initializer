package platform

import (
	"fmt"
	"securebanking-test-data-initializer/pkg/common"
	"time"

	"go.uber.org/zap"
)

//IsValidX509 - Pings the am health endpoint; an invalid certificate will return an error.
//  returns false if no valid cert present after 10 minutes. This is a naive implementation
//  and assumes the error thrown is related to an invalid SSL
func IsValidX509() bool {
	url := fmt.Sprintf("https://%s/am/json/health/live", common.Config.Hosts.IdentityPlatformFQDN)

	for i := 0; i < 60; i++ {
		zap.L().Info("Waiting for valid SSL certificate")
		_, err := restClient.R().
			Get(url)
		if err == nil {
			zap.L().Info("Got valid SSL cert")
			return true
		}
		zap.L().Info(err.Error())
		time.Sleep(10 * time.Second)
	}
	return false
}
