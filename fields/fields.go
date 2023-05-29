package fields

import (
	"context"
	"strings"

	"golang.org/x/exp/slices"
)

type ctxKey uint8

const (
	fieldsCtxKey ctxKey = iota
)

type (
	Field struct {
		Key   string
		Value any
	}
)

func WithFields(ctx context.Context, fields ...Field) context.Context {
	if len(fields) == 0 {
		return ctx
	}

	fs := FromCtx(ctx)
	fs = append(fs, fields...)
	slices.SortFunc(fs, func(a, b Field) bool {
		return strings.Compare(a.Key, b.Key) < 0
	})
	fs = slices.CompactFunc(fs, func(f1, f2 Field) bool {
		return f1.Key == f2.Key
	})

	return context.WithValue(ctx, fieldsCtxKey, fs)
}

func FromCtx(ctx context.Context) []Field {
	fromCtx, _ := ctx.Value(fieldsCtxKey).([]Field)
	return fromCtx
}
