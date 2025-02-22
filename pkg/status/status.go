package status

import (
	"GoGateway/pkg/consts/codes"
	oriErr "errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"net/http"
)

type ErrorWithStatus struct {
	StatusCode int
	Message    error
}

func (e ErrorWithStatus) Error() string {
	return e.Message.Error()
}

func Errorf(code int, format string, a ...interface{}) error {
	return ErrorWithStatus{
		StatusCode: code,
		Message:    errors.Errorf(format, a...),
	}
}

func (e ErrorWithStatus) getErrorStatusCode() int {
	return e.StatusCode
}

type ErrResponse struct {
	Message string `json:"message"`
}

func ErrToHttpResponse(c *app.RequestContext, err error) {
	var v ErrorWithStatus
	if oriErr.As(err, &v) {
		var httpCode int
		var message string

		switch v.getErrorStatusCode() {
		case codes.InternalError:
			hlog.Errorf("Internal Server Error: %v", v.Message)
			httpCode = http.StatusInternalServerError
		case codes.NotFound:
			httpCode = http.StatusNotFound
		case codes.AlreadyExists:
			httpCode = http.StatusConflict
		case codes.InvalidParams:
			httpCode = http.StatusBadRequest
		case codes.Forbidden:
			httpCode = http.StatusForbidden
		case codes.MethodNotAllowed:
			httpCode = http.StatusMethodNotAllowed
		case codes.Unauthorized:
			httpCode = http.StatusUnauthorized
		default:
			httpCode = http.StatusInternalServerError
		}

		message = err.Error()

		if httpCode == http.StatusInternalServerError {
			message = "服务器内部错误"
		}

		c.JSON(httpCode, ErrResponse{Message: message})
		return
	}

	hlog.Error(err.Error())
	c.JSON(http.StatusInternalServerError, ErrResponse{Message: "服务器内部错误"})
}

func NewErrorResponse(message string) *ErrResponse {
	return &ErrResponse{
		Message: message,
	}
}
