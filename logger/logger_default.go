package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type defaultLogger struct {
	maskEnable bool
	zapLogger  *zap.Logger
	level      Level
}

func NewLogger(opt Option) Logger {
	logger := &defaultLogger{}

	if opt.IsEnable {
		logger.zapLogger = NewZapLogger(logger.level, zapcore.AddSync(os.Stdout))
	} else {
		logger.zapLogger = NewZapLogger(logger.level, nil)
	}

	if len(opt.MaskingFields) > 0 {
		logger.maskEnable = true
		maskingFields = make(map[string]bool)

		for _, field := range opt.MaskingFields {
			maskingFields[field] = true
		}
	}

	maskType = FULLY_MASKED
	if opt.MaskingType == PARTIALLY_MASKED {
		maskType = PARTIALLY_MASKED
	}

	if len(opt.AdditionalSkippedContentType) > 0 {
		SkippedContentType = append(SkippedContentType, opt.AdditionalSkippedContentType...)
	}

	Log = logger

	return logger
}

func (d *defaultLogger) Debug(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{
		zap.String("level", "debug"),
	}

	traceContextFields := TraceContext(ctx)
	fields := formatToField(details...)
	zapLogs = append(zapLogs, formatLogs(ctx, message, d.maskEnable, fields...)...)
	d.zapLogger.With(traceContextFields...).Debug(message, zapLogs...)
}

func (d *defaultLogger) Info(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{
		zap.String("level", "info"),
	}

	traceContextFields := TraceContext(ctx)

	fields := formatToField(details...)
	zapLogs = append(zapLogs, formatLogs(ctx, message, d.maskEnable, fields...)...)
	d.zapLogger.With(traceContextFields...).Info(message, zapLogs...)
}

func (d *defaultLogger) Warn(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{
		zap.String("level", "warn"),
	}

	traceContextFields := TraceContext(ctx)

	fields := formatToField(details...)
	zapLogs = append(zapLogs, formatLogs(ctx, message, d.maskEnable, fields...)...)
	d.zapLogger.With(traceContextFields...).Warn(message, zapLogs...)
}

func (d *defaultLogger) Error(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{
		zap.String("level", "error"),
	}

	traceContextFields := TraceContext(ctx)

	fields := formatToField(details...)
	zapLogs = append(zapLogs, formatLogs(ctx, message, d.maskEnable, fields...)...)
	d.zapLogger.With(traceContextFields...).Error(message, zapLogs...)
}

func (d *defaultLogger) Fatal(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{
		zap.String("level", "fatal"),
	}

	traceContextFields := TraceContext(ctx)

	fields := formatToField(details...)
	zapLogs = append(zapLogs, formatLogs(ctx, message, d.maskEnable, fields...)...)
	d.zapLogger.With(traceContextFields...).Fatal(message, zapLogs...)
}

func (d *defaultLogger) Panic(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{
		zap.String("level", "panic"),
	}

	traceContextFields := TraceContext(ctx)

	fields := formatToField(details...)
	zapLogs = append(zapLogs, formatLogs(ctx, message, d.maskEnable, fields...)...)
	d.zapLogger.With(traceContextFields...).Panic(message, zapLogs...)
}

func (d *defaultLogger) TDR(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{
		zap.String("level", "TDR"),
	}

	traceContextFields := TraceContext(ctx)

	fields := formatToField(details...)
	zapLogs = append(zapLogs, formatLogsTDR(ctx, message, d.maskEnable, fields...)...)
	d.zapLogger.With(traceContextFields...).Info(message, zapLogs...)
}

func formatLogs(ctx context.Context, msg string, mask bool, fields ...Field) (logRecord []zap.Field) {
	ctxVal := ExtractCtx(ctx)

	// add global value from context that must be exist on all logs!
	logRecord = append(logRecord, zap.String("message", msg))

	logRecord = append(logRecord, zap.String("_app_name", ctxVal.ServiceName))
	logRecord = append(logRecord, zap.String("_app_version", ctxVal.ServiceVersion))
	logRecord = append(logRecord, zap.Int("_app_port", ctxVal.ServicePort))
	logRecord = append(logRecord, zap.String("_app_tag", ctxVal.Tag))
	logRecord = append(logRecord, zap.String("_app_method", ctxVal.ReqMethod))
	logRecord = append(logRecord, zap.String("_app_uri", ctxVal.ReqURI))

	// add additional data that available across all log, such as user_id
	if ctxVal.AdditionalData != nil {
		logRecord = append(logRecord, zap.Any("_app_data", ctxVal.AdditionalData))
	}

	for _, field := range fields {
		logRecord = append(logRecord, formatLog(field.Key, field.Val, mask))
	}

	return
}

func formatLogsTDR(ctx context.Context, msg string, mask bool, fields ...Field) []zap.Field {
	ctxVal := ExtractCtx(ctx)

	logRecord := []zap.Field{
		zap.String("message", msg),
		zap.String("_app_name", ctxVal.ServiceName),
		zap.String("_app_version", ctxVal.ServiceVersion),
		zap.Int("_app_port", ctxVal.ServicePort),
		zap.String("_app_tag", ctxVal.Tag),
		zap.String("_app_method", ctxVal.ReqMethod),
		zap.String("_app_uri", ctxVal.ReqURI),
		zap.Int("_app_response_code", ctxVal.RespCode),
		zap.String("_app_exec_time", fmt.Sprintf("%sms", ctxVal.RespTime)),
	}

	logRecord = append(logRecord, formatLog("_app_request_header", ctxVal.ReqHeader, mask))
	logRecord = append(logRecord, formatLog("_app_request", ctxVal.ReqBody, mask))
	logRecord = append(logRecord, formatLog("_app_response", ctxVal.RespBody, mask))

	if ctxVal.AdditionalData != nil {
		logRecord = append(logRecord, zap.Any("_app_data", ctxVal.AdditionalData))
	}

	if ctxVal.Error != "" {
		logRecord = append(logRecord, zap.Any("_app_error", ctxVal.Error))
	}

	for _, field := range fields {
		logRecord = append(logRecord, formatLog(field.Key, field.Val, mask))
	}

	return logRecord
}

func formatToField(details ...interface{}) (logRecord []Field) {
	for index, msg := range details {
		logRecord = append(logRecord, Field{
			Key: "_message_" + cast.ToString(index),
			Val: msg,
		},
		)
	}

	return
}

func formatLog(key string, msg interface{}, mask bool) (logRecord zap.Field) {
	if msg == nil {
		logRecord = zap.Any(key, struct{}{})
		return
	}

	// handle proto message
	p, ok := msg.(proto.Message)
	if ok {
		b, _err := json.Marshal(p)
		if _err != nil {
			logRecord = zap.Any(key, p.String())
			return
		}

		var data interface{}
		if _err = json.Unmarshal(b, &data); _err != nil {
			// string cannot be masked, so only try to marshal as json object
			logRecord = zap.Any(key, p.String())
			return
		}

		if !mask {
			logRecord = zap.Any(key, data)
			return
		}

		logRecord = zap.Any(key, MaskData(data, maskingFields))
		return
	}

	// handle string, string is cannot be masked, just write it
	// but if parse as json object sucess, then masked it
	if str, ok := msg.(string); ok {
		var data interface{}
		if _err := json.Unmarshal([]byte(str), &data); _err != nil {
			logRecord = zap.String(key, str)
			return
		}

		if !mask {
			logRecord = zap.Any(key, data)
			return
		}

		logRecord = zap.Any(key, MaskData(data, maskingFields))
		return
	}

	// if masking is disabled then just set as field log
	if !mask {
		logRecord = zap.Any(key, msg)
		return
	}

	logRecord = zap.Any(key, MaskData(msg, maskingFields))
	return
}
