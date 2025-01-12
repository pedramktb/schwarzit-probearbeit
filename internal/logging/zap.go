package logging

import (
	"context"

	"go.uber.org/zap"
)

type ContextKey string

const (
	CtxUserID ContextKey = "user.id"
)

var (
	ctxKeys = []ContextKey{
		CtxUserID,
	}
)

func FromContext(cx context.Context) *zap.Logger {
	if cx == nil {
		return zap.L()
	}
	return zap.L().With(ToZapFields(cx)...)
}

func ToZapFields(cx context.Context) []zap.Field {
	var fields []zap.Field
	for _, key := range ctxKeys {
		if value := cx.Value(key); value != nil {
			strKey := string(key)

			switch v := value.(type) {
			case string:
				if v != "" {
					fields = append(fields, zap.String(strKey, v))
				}
			case int:
				fields = append(fields, zap.Int(strKey, v))
			case int8:
				fields = append(fields, zap.Int8(strKey, v))
			case int16:
				fields = append(fields, zap.Int16(strKey, v))
			case int32:
				fields = append(fields, zap.Int32(strKey, v))
			case int64:
				fields = append(fields, zap.Int64(strKey, v))
			case bool:
				fields = append(fields, zap.Bool(strKey, v))
			case float32:
				fields = append(fields, zap.Float32(strKey, v))
			case float64:
				fields = append(fields, zap.Float64(strKey, v))
			default:
				fields = append(fields, zap.Any(strKey, value))
			}
		}
	}
	return fields
}
