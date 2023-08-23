package errors

import (
	"context"
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

	oldFields := FromCtx(ctx)
	oldFieldIdxByKey := make(map[string]int, len(oldFields))
	for idx, f := range oldFields {
		oldFieldIdxByKey[f.Key] = idx
	}

	var (
		wasCopied bool
		result    = oldFields
	)

	for _, newField := range fields {
		oldFieldIdx, ok := oldFieldIdxByKey[newField.Key]
		if !ok {
			result = append(result, newField)
			continue
		}

		oldField := oldFields[oldFieldIdx]
		// skip copying if field not changed
		if oldField == newField {
			continue
		}

		if !wasCopied {
			tmp := make([]Field, len(result))
			copy(tmp, result)
			result = tmp
			wasCopied = true
		}

		result[oldFieldIdx] = newField
	}

	return context.WithValue(ctx, fieldsCtxKey, result)
}

// FromCtx extract field from context.
func FromCtx(ctx context.Context) []Field {
	fromCtx, _ := ctx.Value(fieldsCtxKey).([]Field)
	return fromCtx
}
