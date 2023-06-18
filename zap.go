package errors

import (
	"errors"

	"go.uber.org/zap"
)

func ZapFields(err error) (fs []zap.Field) {
	var e *ExtendedError
	fs = append(fs, zap.Error(err))
	if !errors.As(err, &e) {
		return fs
	}

	if stack := e.Stacktrace(); stack != "" {
		fs = append(fs, zap.String("stacktrace", stack))
	}
	for _, f := range e.Fields() {
		fs = append(fs, zap.Any(f.Key, f.Value))
	}

	return fs
}
