package common

import (
	"go.uber.org/zap"
)

// RestError for use with goresty SetError client method
//    Allow us to marshal client and server error responses
type RestError struct {
	Message string
	Code    int
	Detail  interface{}
}

// RaiseForStatus will exit if go resty returns an error in STRICT mode,
//    Be it client error, server error or other. Turning off
//    STRICT mode will simply warn of client/server errors.
func RaiseForStatus(err error, restError interface{}, status int) {
	if err != nil {
		zap.S().Fatalw("Go rest has thrown an error when attempting to send", "error", err, "httpStatus", status)
	}

	if restError != nil {

		strict := Config.Environment.Strict
		if strict {
			zap.S().Fatalw("Go rest has sent the request but the status is > 399", "error", restError, "httpStatus", status)
		}
		zap.S().Warnw("Go rest has sent the request but the status is > 399", "error", restError, "httpStatus", status)
	}
}
