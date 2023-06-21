package errors

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleError(t *testing.T) {
	err := NewCtx(context.Background(), "test msg with args %d %s", 1, "2")
	assert.Error(t, err)
	assert.Equal(t, "test msg with args 1 2", err.Error())
}

func TestStackTraceDiscovery(t *testing.T) {
	err := WrapErr(
		WrapExtErr(
			WrapErr(
				Call(),
			),
		),
	)

	assert.Error(t, err)
	var ext *ExtendedError
	assert.True(t, As(err, &ext))
	assert.Equal(t, "github.com/popoffvg/dev-errors.Call\n\t/home/popoffvg/Documents/git/dev-errors/error_test.go:50\ngithub.com/popoffvg/dev-errors.TestStackTraceDiscovery\n\t/home/popoffvg/Documents/git/dev-errors/error_test.go:22\ntesting.tRunner\n\t/usr/lib/go-1.20/src/testing/testing.go:1576", ext.Stacktrace())
}

func TestStackTraceUnion(t *testing.T) {
	err := MultiError()
	assert.Error(t, err)
	var ext *ExtendedError
	assert.True(t, As(err, &ext))
	assert.Equal(t, "github.com/popoffvg/dev-errors.Call\n\t/home/popoffvg/Documents/git/dev-errors/error_test.go:50\ngithub.com/popoffvg/dev-errors.MultiError\n\t/home/popoffvg/Documents/git/dev-errors/error_test.go:62\ngithub.com/popoffvg/dev-errors.TestStackTraceUnion\n\t/home/popoffvg/Documents/git/dev-errors/error_test.go:34\ntesting.tRunner\n\t/usr/lib/go-1.20/src/testing/testing.go:1576\n\n--------------------\n\ngithub.com/popoffvg/dev-errors.Call\n\t/home/popoffvg/Documents/git/dev-errors/error_test.go:50\ngithub.com/popoffvg/dev-errors.MultiError\n\t/home/popoffvg/Documents/git/dev-errors/error_test.go:62\ngithub.com/popoffvg/dev-errors.TestStackTraceUnion\n\t/home/popoffvg/Documents/git/dev-errors/error_test.go:34\ntesting.tRunner\n\t/usr/lib/go-1.20/src/testing/testing.go:1576", ext.Stacktrace())
}

func TestWVerb(t *testing.T) {
	oldErr := errors.New("old error")
	err := New("test: %w", oldErr)
	assert.True(t, Is(err, oldErr))

	assert.Equal(t, "test: old error", err.Error())
}

func Call() error {
	return New("from Call")
}

func WrapErr(err error) error {
	return fmt.Errorf("wrap %w", err)
}

func WrapExtErr(err error) error {
	return New("%s", err)
}

func MultiError() error {
	return New("%s", fmt.Errorf("%w %w", Call(), WrapErr(Call())))
}
