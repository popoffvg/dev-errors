package errors_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/popoffvg/dev-errors/errors"
	"github.com/stretchr/testify/assert"
)

func TestSimpleError(t *testing.T) {
	err := errors.NewCtx(context.Background(), "test msg with args %d %s", 1, "2")
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
	var ext *errors.ExtendedError
	assert.True(t, errors.As(err, &ext))
	assert.Equal(t, "github.com/popoffvg/dev-errors/errors_test.Call\n\t/home/popoffvg/Documents/git/dev-errors/errors/error_test.go:42\ngithub.com/popoffvg/dev-errors/errors_test.TestStackTraceDiscovery\n\t/home/popoffvg/Documents/git/dev-errors/errors/error_test.go:22\ntesting.tRunner\n\t/usr/lib/go-1.20/src/testing/testing.go:1576", ext.Stacktrace())
}

func TestStackTraceUnion(t *testing.T) {
	err := MultiError()
	assert.Error(t, err)
	var ext *errors.ExtendedError
	assert.True(t, errors.As(err, &ext))
	assert.Equal(t, "github.com/popoffvg/dev-errors/errors_test.Call\n\t/home/popoffvg/Documents/git/dev-errors/errors/error_test.go:42\ngithub.com/popoffvg/dev-errors/errors_test.MultiError\n\t/home/popoffvg/Documents/git/dev-errors/errors/error_test.go:54\ngithub.com/popoffvg/dev-errors/errors_test.TestStackTraceUnion\n\t/home/popoffvg/Documents/git/dev-errors/errors/error_test.go:34\ntesting.tRunner\n\t/usr/lib/go-1.20/src/testing/testing.go:1576\n\n--------------------\n\ngithub.com/popoffvg/dev-errors/errors_test.Call\n\t/home/popoffvg/Documents/git/dev-errors/errors/error_test.go:42\ngithub.com/popoffvg/dev-errors/errors_test.MultiError\n\t/home/popoffvg/Documents/git/dev-errors/errors/error_test.go:54\ngithub.com/popoffvg/dev-errors/errors_test.TestStackTraceUnion\n\t/home/popoffvg/Documents/git/dev-errors/errors/error_test.go:34\ntesting.tRunner\n\t/usr/lib/go-1.20/src/testing/testing.go:1576", ext.Stacktrace())
}

func Call() error {
	return errors.New("from Call")
}

func WrapErr(err error) error {
	return fmt.Errorf("wrap %w", err)
}

func WrapExtErr(err error) error {
	return errors.New("%s", err)
}

func MultiError() error {
	return errors.New("%s", fmt.Errorf("%w %w", Call(), WrapErr(Call())))
}
