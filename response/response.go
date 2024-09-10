package response

import (
	"fmt"
	"go-common/constant"
	"go-common/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MetaTpl struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalData int `json:"total_data"`
}

type BasePayload struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

func WriteError(c *gin.Context, err error) {

	c.Set(constant.ErrorMessageKey, err.Error())

	httpStatusCode := http.StatusInternalServerError
	payload := BasePayload{
		Error:   errors.ErrorCodeGeneralError,
		Success: false,
		Message: fmt.Sprintf("fatal error: %s", err.Error()),
	}

	if err, ok := err.(*errors.Err); ok {
		payload.Message = err.Error()
		payload.Error = err.GetErrorCode()
		httpStatusCode = err.GetHttpStatusCode()
	}

	c.JSON(httpStatusCode, payload)
}
