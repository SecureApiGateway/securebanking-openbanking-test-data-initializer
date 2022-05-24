package platform

import (
	"github.com/go-resty/resty/v2"
	"securebanking-test-data-initializer/pkg/common"
)

var restClient = resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).SetError(common.RestError{})
