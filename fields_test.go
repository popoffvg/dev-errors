package errors

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFields(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
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

		ctx := WithFields(context.Background(), fields...)

		expected := fields

		actual := FromCtx(ctx)
		require.NotEmpty(t, actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("with added new field", func(t *testing.T) {
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

		ctx := WithFields(context.Background(), fields...)
		require.NotNil(t, ctx)

		newFields := []Field{
			{
				Key:   "key-1",
				Value: 10,
			},
			{
				Key:   "key-2",
				Value: 20,
			},
			{
				Key:   "key-3",
				Value: 30,
			},
		}

		ctx = WithFields(ctx, newFields...)
		require.NotNil(t, ctx)

		expected := []Field{
			{
				Key:   "key-2",
				Value: 20,
			},
			{
				Key:   "key-1",
				Value: 10,
			},
			{
				Key:   "key-3",
				Value: 30,
			},
		}

		actual := FromCtx(ctx)
		require.NotEmpty(t, ctx)
		assert.Equal(t, expected, actual)
	})

	t.Run("with goroutines", func(t *testing.T) {
		fields := []Field{
			{
				Key:   "key-1",
				Value: 1,
			},
			{
				Key:   "key-2",
				Value: 2,
			},
		}

		newFields := []Field{
			{
				Key:   "key-2",
				Value: 10,
			},
			{
				Key:   "key-1",
				Value: 3,
			},
		}

		ctx := WithFields(context.Background(), fields...)
		require.NotNil(t, ctx)

		ch := make(chan struct{})
		wg := sync.WaitGroup{}

		var fieldsInFirstGoroutine []Field
		wg.Add(1)
		go func() {
			defer wg.Done()
			fieldsInFirstGoroutine = FromCtx(ctx)
			ch <- struct{}{}
		}()

		var fieldsInSecondGoroutine []Field
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ch
			ctx = WithFields(ctx, newFields...)
			fieldsInSecondGoroutine = FromCtx(ctx)
		}()

		wg.Wait()

		require.NotNil(t, fieldsInFirstGoroutine)
		require.NotNil(t, fieldsInSecondGoroutine)
		assert.NotEqual(t, fieldsInFirstGoroutine, fieldsInSecondGoroutine)
	})
}
