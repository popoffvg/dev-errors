package errors

import (
	"context"

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

// WithFields added field to context.
func WithFields(ctx context.Context, fields ...Field) context.Context {
	if len(fields) == 0 {
		return ctx
	}

	fs := FromCtx(ctx)

	var (
		wasCopied bool
		result    = fs
	)
	for _, f := range fields {
		j := slices.IndexFunc(result, func(v Field) bool {
			return v.Key == f.Key
		})
		if j != -1 {
			// skip copying if value not changed
			if result[j].Value == f.Value {
				continue
			}

			if !wasCopied {
				tmp := make([]Field, len(result))
				copy(tmp, result)
				result = tmp
			}
			result[j] = f
			continue
		}

		result = append(result, f)
	}

	return context.WithValue(ctx, fieldsCtxKey, result)
}

// FromCtx extract field from context.
func FromCtx(ctx context.Context) []Field {
	fromCtx, _ := ctx.Value(fieldsCtxKey).([]Field)
	return fromCtx
}
