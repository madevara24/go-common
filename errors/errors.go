package errors

import "net/http"

const (
	ErrorCodeGeneralError = "-1"
	ErrorInvalidParameter = "PAYLOAD-99"
)

var (
	ErrJSONParse        = NewErr(http.StatusBadRequest, "PAYLOAD-01", "invalid json payload")
	ErrJWTClaimsParse   = NewErr(http.StatusBadRequest, "PAYLOAD-02", "unable read user claims")
	ErrUserTokenMissing = NewErr(http.StatusBadRequest, "PAYLOAD-03", "unable retrieve user token")
	ErrUnauthorized     = NewErr(http.StatusUnauthorized, "PAYLOAD-04", "token invalid")

	ErrTimeFormat = NewErr(http.StatusInternalServerError, "DATA-01", "invalid time format")
)

type Err struct {
	errorCode      string
	message        string
	httpStatusCode int
}

func (e *Err) Error() string {
	return e.message
}

func (e *Err) GetHttpStatusCode() int {
	return e.httpStatusCode
}

func (e *Err) GetErrorCode() string {
	return e.errorCode
}

func NewErr(httpStatusCode int, errorCode, message string) *Err {
	return &Err{
		httpStatusCode: httpStatusCode,
		errorCode:      errorCode,
		message:        message,
	}
}
