package httprest

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"log"
	"net/http"
	"securebanking-test-data-initializer/pkg/common"
)

var Client RestReaderWriter

func InitRestReaderWriter(cookie *http.Cookie, authCode string) {
	Client = &RestClient{
		Resty:                resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).SetError(common.RestError{}),
		Cookie:               cookie,
		AuthCode:             authCode,
		IdentityPlatformFQDN: common.Config.Hosts.Scheme + "://" + common.Config.Hosts.IdentityPlatformFQDN,
	}
}

func (r *RestClient) Get(path string, headers map[string]string) ([]byte, int) {
	route := r.IdentityPlatformFQDN + path
	resp, err := r.request(headers).
		Get(route)
	log.Println("Route:", route, resp.Status())
	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	return resp.Body(), resp.StatusCode()
}

func (r *RestClient) GetRS(path string, headers map[string]string) ([]byte, int) {
	resp, err := r.request(headers).
		Get(path)
	log.Println("Route:", path, resp.Status())
	if err != nil {
		zap.S().Infow("Error request", "path", path, "error", err, "status", resp.StatusCode())
	}
	//common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	return resp.Body(), resp.StatusCode()
}

func (r *RestClient) request(headers map[string]string) *resty.Request {
	return r.Resty.R().
		SetHeaders(headers).
		SetCookie(r.Cookie).
		SetAuthToken(r.AuthCode)
}

func (r *RestClient) Post(path string, ob interface{}, headers map[string]string) ([]byte, int) {
	resp, err := r.request(headers).
		SetBody(ob).
		SetContentLength(true).
		Post(r.IdentityPlatformFQDN + path)
	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	return resp.Body(), resp.StatusCode()
}

func (r *RestClient) PostRS(path string, headers map[string]string) int {
	resp, err := r.request(headers).
		SetContentLength(true).
		Post(path)
	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	return resp.StatusCode()
}

func (r *RestClient) Patch(path string, ob interface{}, headers map[string]string) int {
	resp, err := r.request(headers).
		SetBody(ob).
		Patch(r.IdentityPlatformFQDN + path)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	return resp.StatusCode()
}

func (r *RestClient) Put(path string, ob interface{}, headers map[string]string) int {
	resp, err := r.request(headers).
		SetBody(ob).
		SetContentLength(true).
		Put(r.IdentityPlatformFQDN + path)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	return resp.StatusCode()
}
