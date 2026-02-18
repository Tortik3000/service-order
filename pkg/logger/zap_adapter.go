package logger

import (
	"go.uber.org/zap"
)

type zapAdapter struct {
	logger *zap.Logger
}

var _ Logger = (*zapAdapter)(nil)

func NewZap(
	logger *zap.Logger,
) *zapAdapter {
	return &zapAdapter{
		logger: logger,
	}
}

func (z *zapAdapter) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, convertFields(fields)...)
}

func (z *zapAdapter) Info(msg string, fields ...Field) {
	z.logger.Info(msg, convertFields(fields)...)
}

func (z *zapAdapter) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, convertFields(fields)...)
}

func (z *zapAdapter) Error(msg string, fields ...Field) {
	z.logger.Error(msg, convertFields(fields)...)
}

func (z *zapAdapter) Fatal(msg string, fields ...Field) {
	z.logger.Fatal(msg, convertFields(fields)...)
}

func (z *zapAdapter) With(fields ...Field) Logger {
	return NewZap(z.logger.With(convertFields(fields)...))
}

func (z *zapAdapter) Sync() error {
	return z.logger.Sync()
}

func convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
