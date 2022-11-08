package log

import "go.uber.org/zap"

func MapToFields(m map[string]string) []zap.Field {
	fields := make([]zap.Field, 0, len(m))
	for k, v := range m {
		fields = append(fields, zap.String(k, v))
	}
	return fields
}
