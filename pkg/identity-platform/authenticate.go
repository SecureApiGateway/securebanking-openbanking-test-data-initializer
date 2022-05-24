package platform

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"securebanking-test-data-initializer/pkg/common"
	"securebanking-test-data-initializer/pkg/types"
)

func GetCookieNameFromAm() string {
	zap.L().Info("Getting Cookie name from Identity Platform")
	path := fmt.Sprintf("%s://%s/am/json/serverinfo/*", common.Config.Hosts.Scheme, common.Config.Hosts.IdentityPlatformFQDN)
	zap.S().Infow("Getting Cookie name from Identity Platform", "path", path)
	result := &types.ServerInfo{}
	resp, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetResult(result).
		Get(path)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	cookieName := result.CookieName

	zap.S().Infow("Got cookie from am",
		zap.String("cookieName", cookieName))
	return cookieName
}

// FromUserSession - get a session token from AM for authentication
//    returns the Session object with embedded session cookie
func FromUserSession(cookieName string) *common.Session {
	zap.L().Info("Getting an admin session from Identity Platform")

	path := fmt.Sprintf("https://%s/am/json/realms/root/authenticate?authIndexType=service&authIndexValue=ldapService", common.Config.Hosts.IdentityPlatformFQDN)

	zap.S().Infow("Path to authenticate the user", "path", path)

	resp, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Accept-API-Version", "resource=2.0, protocol=1.0").
		SetHeader("X-OpenAM-Username", common.Config.Users.CdmAdminUsername).
		SetHeader("X-OpenAM-Password", common.Config.Users.CdmAdminPassword).
		Post(path)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	var cookieValue = ""
	for _, cookie := range resp.Cookies() {
		zap.S().Infow("Cookie found", "cookie", cookie)
		if cookie.Name == cookieName {
			cookieValue = cookie.Value
		}
	}
	if cookieValue == "" {
		zap.S().Fatalw("Cannot find cookie",
			"statusCode", resp.StatusCode(),
			"cookieName", cookieName,
			"advice", `Calling this method twice might result in this error,
				 AM will not react well to successive session requests`,
			"error", resp.Error())
	}
	c := &http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Domain:   common.Config.Hosts.IdentityPlatformFQDN,
	}
	s := &common.Session{}
	s.Cookie = c
	zap.S().Infow("New Identity Platform session created", "cookie", s.Cookie)
	return s
}
