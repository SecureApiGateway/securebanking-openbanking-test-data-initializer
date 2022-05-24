package httprest

import (
	"github.com/go-resty/resty/v2"
	"net/http"
)

type RestReaderWriter interface {
	Get(string, map[string]string) ([]byte, int)
	GetRS(string, map[string]string) ([]byte, int)
	Patch(string, interface{}, map[string]string) int
	Post(string, interface{}, map[string]string) ([]byte, int)
	PostRS(string, map[string]string) int
	Put(string, interface{}, map[string]string) int
}

type RestClient struct {
	Resty                *resty.Client
	Cookie               *http.Cookie
	AuthCode             string
	IdentityPlatformFQDN string
}
