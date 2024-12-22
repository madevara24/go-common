package logger

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
)

var SkippedContentType = []string{
	"application/tar+gzip",
	"application/octet-stream",
	"application/msword",
	"application/zip",
	"application/pdf",
	"image/gif",
	"image/jpeg",
	"image/png",
	"image/tiff",
	"application/vnd",
	"multipart/form-data",
}

const (
	FULLY_MASKED     = "FULL"
	PARTIALLY_MASKED = "PARTIAL"
)

func IsSkipLog(contentType string) bool {
	for _, list := range SkippedContentType {
		if strings.Contains(contentType, list) {
			return true
		}
	}

	return false
}

func MaskData(data interface{}, sensitiveFields map[string]bool) interface{} {
	switch value := data.(type) {
	case map[string]interface{}:
		maskMapFields(value, sensitiveFields)
	case []interface{}:
		maskSliceElements(value, sensitiveFields)
	case interface{}:
		dataMap := convertToMap(value)
		if dataMap != nil {
			MaskData(dataMap, sensitiveFields)
		}
		// Do nothing for non-map and non-slice values
	}
	return data
}

func maskMapFields(dataMap map[string]interface{}, sensitiveFields map[string]bool) {
	for k, v := range dataMap {
		if sensitiveFields[k] {
			if maskType == PARTIALLY_MASKED {
				if str, ok := v.(string); ok && len(str) >= 2 {
					dataMap[k] = "***" + str[len(str)-2:]
				} else {
					dataMap[k] = "*****"
				}
			} else {
				dataMap[k] = "*****"
			}
		} else {
			MaskData(v, sensitiveFields)
		}
	}
}

func maskSliceElements(slice []interface{}, sensitiveFields map[string]bool) {
	for _, v := range slice {
		MaskData(v, sensitiveFields)
	}
}

func convertToMap(value interface{}) map[string]interface{} {
	by, err := json.Marshal(value)
	if err != nil {
		return nil
	}

	tempData := map[string]interface{}{}
	if err := json.Unmarshal(by, &tempData); err != nil {
		return nil
	}

	return tempData
}

func TraceContext(ctx context.Context) []zapcore.Field {
	const (
		FieldKeyTraceID = "trace.id"
		FieldKeySpanID  = "span.id"
	)

	spanContext := trace.SpanFromContext(ctx).SpanContext()
	fields := []zapcore.Field{
		zap.Stringer(FieldKeyTraceID, spanContext.TraceID()),
		zap.Stringer(FieldKeySpanID, spanContext.SpanID()),
	}

	return fields
}

func IsMultipartForm(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "multipart/form-data")
}

func ParseMultipartFields(r *http.Request) (map[string]string, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	for key, values := range r.MultipartForm.Value {
		if len(values) > 0 {
			fields[key] = values[0]
		} else {
			fields[key] = ""
		}
	}

	return fields, nil
}
