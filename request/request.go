package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/madevara24/go-common/errors"
	"github.com/madevara24/go-common/response"
)

func UnmarshalJSON(c *gin.Context, dest interface{}) error {
	if c.Request.Body == nil {
		return nil
	}

	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.WriteError(c, errors.ErrJSONParse)
		return err
	}

	if err := json.Unmarshal(bodyBytes, dest); err != nil {
		response.WriteError(c, errors.ErrJSONParse)
		return err
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return nil
}
