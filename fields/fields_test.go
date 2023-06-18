package fields

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFields(t *testing.T) {
	fields := []Field{
		{
			Key:   "key-2",
			Value: 2,
		},
		{
			Key:   "key-1",
			Value: 1,
		},
	}

	ctx := context.Background()
	ctx = WithFields(ctx, fields...)
	assert.Equal(t, fields, FromCtx(ctx))

	ctx = WithFields(ctx, Field{"key-3", 3}, Field{"key-1", 5})
	fields[1].Value = 5
	assert.ElementsMatch(t, append(fields, Field{"key-3", 3}), FromCtx(ctx))
}
