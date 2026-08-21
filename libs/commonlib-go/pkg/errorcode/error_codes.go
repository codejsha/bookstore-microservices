package errorcode

import (
	"net/http"
)

type errorLevel string

const (
	errorLevelUnknown errorLevel = "UNKNOWN"
	errorLevelDebug   errorLevel = "DEBUG"
	errorLevelInfo    errorLevel = "INFO"
	errorLevelWarning errorLevel = "WARNING"
	errorLevelError   errorLevel = "ERROR"
)

type ErrorCode struct {
	Code       int        `json:"code"`
	ErrorLevel errorLevel `json:"error_level"`
	Message    string     `json:"message"`
	StatusCode int        `json:"-"`
}

var (
	UnknownError = ErrorCode{Code: 10001, ErrorLevel: errorLevelError, Message: "unknown error", StatusCode: http.StatusInternalServerError}
)

var (
	InvalidRequestPayload = ErrorCode{Code: 20001, ErrorLevel: errorLevelError, Message: "invalid request payload", StatusCode: http.StatusBadRequest}
	InvalidRequestParams  = ErrorCode{Code: 20002, ErrorLevel: errorLevelError, Message: "invalid request params", StatusCode: http.StatusBadRequest}
)
