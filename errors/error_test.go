package errors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleError(t *testing.T) {
	err := NewCtx(context.Background(), "test msg with args %d %s", 1, "2")
	assert.Error(t, err)
	assert.Equal(t, "", err.Error())
}
