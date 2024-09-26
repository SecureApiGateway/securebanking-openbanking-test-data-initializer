package platform

import (
	"encoding/json"
	"fmt"
        "github.com/go-jose/go-jose/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"securebanking-test-data-initializer/pkg/common"
	"securebanking-test-data-initializer/pkg/types"
	"time"
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
//
//	returns the Session object with embedded session cookie
func FromUserSession(cookieName string) *common.Session {
	zap.L().Info("Getting an admin session from Identity Platform")

	path := ""
	path = fmt.Sprintf("https://%s/am/json/realms/root/authenticate?authIndexType=service&authIndexValue=ldapService", common.Config.Hosts.IdentityPlatformFQDN)

	zap.S().Infow("Path to authenticate the user", "path", path)

	resp, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Accept-API-Version", "resource=2.0, protocol=1.0").
		SetHeader("X-OpenAM-Username", common.Config.Users.CDKPlatformAdminUsername).
		SetHeader("X-OpenAM-Password", common.Config.Users.CDKPlatformAdminPassword).
		Post(path)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())
	zap.S().Infof("Got response code %v from %v", resp.StatusCode(), path)

	var cookieValue = ""
	for _, cookie := range resp.Cookies() {
		zap.S().Infow("Cookies found", "cookie", cookie)
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

func GetServiceAccountToken() string {

	zap.S().Infof("Getting token with service account")

	serviceAccountId := common.Config.Users.FIDCPlatformServiceAccountId
	serviceAccountKey := common.Config.Users.FIDCPlatformServiceAccountKey

	if serviceAccountId == "" || serviceAccountKey == "" {
		zap.S().Fatalw("Service account details not set.")
	} else if serviceAccountId == "replaceme" || serviceAccountKey == "replaceme" {
		zap.S().Fatalw("Service account details have not been overwritten from default values.")
	}

	serviceAccountKeyJWK := jose.JSONWebKey{}
	serviceAccountKeyJWK.UnmarshalJSON([]byte(serviceAccountKey))

	// Instantiate a signer using RS256 with the given private key.
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: "RS256", Key: serviceAccountKeyJWK}, nil)
	if err != nil {
		panic(err)
	}

	tokenEndpoint := fmt.Sprintf("https://%s/am/oauth2/access_token", common.Config.Hosts.IdentityPlatformFQDN)
	const JWT_VALIDITY_SECONDS = 180

	zap.S().Infof("Using token endpoint: %v", tokenEndpoint)

	payload := types.JWTPayload{
		ISS: serviceAccountId,
		SUB: serviceAccountId,
		AUD: tokenEndpoint,
		JTI: uuid.NewString(),
		EXP: time.Now().Unix() + JWT_VALIDITY_SECONDS,
	}

	payloadByte, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	jws, err := signer.Sign(payloadByte)
	if err != nil {
		panic(err)
	}

	jwt, err := jws.CompactSerialize()
	if err != nil {
		panic(err)
	}

	zap.S().Infof("Generated JWT for token request: %v", jwt)

	resp, err := restClient.R().
		SetFormData(map[string]string{
			"grant_type": "urn:ietf:params:oauth:grant-type:jwt-bearer",
			"client_id":  "service-account",
			"scope":      "fr:idm:* fr:am:* fr:idc:esv:*",
			"assertion":  jwt,
		}).
		Post(tokenEndpoint)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())
	zap.S().Infof("Got response code %v from %v with body %v", resp.StatusCode(), tokenEndpoint, string(resp.Body()))

	var responseJson map[string]interface{}
	json.Unmarshal(resp.Body(), &responseJson)
	token := responseJson["access_token"]

	zap.S().Infof("Got access token for service account: %v", token)

	return fmt.Sprintf("%v", token)
}
