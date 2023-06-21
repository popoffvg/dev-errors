package errors

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/popoffvg/dev-errors/internal/buffer"
	"github.com/popoffvg/dev-errors/internal/bufferpool"
)

type (
	ExtendedError struct {
		stack *stacktrace

		causes []error
		msg    string
		fields []Field
	}
)

// New create new error.
//
// Stack will capture if captureStack option isn't off.
func New(msg string, args ...any) error {
	return newErr(context.Background(), msg, false, args...)
}

// NewCtx create new error.
//
// Stack will capture if captureStack option isn't off.
// Fields from context will added to error if EnableField option isn't off.
func NewCtx(ctx context.Context, msg string, args ...any) error {
	return newErr(ctx, msg, false, args...)
}

// Wrap create new error and added other error as cause.
func Wrap(err error) error {
	return newErr(context.Background(), "", true, err)
}

// Wrap create new error and added other error as cause
// and fields from context.
func WrapCtx(ctx context.Context, err error) error {
	return newErr(ctx, "", true, err)
}

// Unwrap implement Unwrap interface from "errors" pkg.
func (e *ExtendedError) Unwrap() []error {
	if len(e.causes) == 0 {
		return nil
	}

	return e.causes
}

// Error implement Error interface from "errors" pkg.
func (e *ExtendedError) Error() string {
	return opts.printer.Print(e.msg, e.frames(), e.Fields())
}

// Fields return fields saved from context.
func (e *ExtendedError) Fields() []Field {
	return e.fields
}

// Stacktrace return stack as string.
func (e *ExtendedError) Stacktrace() string {
	var buf = bufferpool.Get()
	writeStack(buf, e.frames())
	return buf.String()
}

func newErr(ctx context.Context, msg string, skipMsg bool, args ...any) error {
	var (
		causes []error

		hasStack       bool
		causeWithStack = new(ExtendedError)
	)

	// supported %w verb: add error to cause and replace to %s verb
	msg = strings.ReplaceAll(msg, "%w", "%s")
	if opts.withStack {
		for i, v := range args {
			if e, ok := v.(error); ok {
				causes = append(causes, e)
				args[i] = e.Error()
				if !hasStack && errors.As(e, &causeWithStack) && causeWithStack.stack != nil {
					hasStack = true
				}
			}
		}
	}

	return applyHook(&ExtendedError{
		stack: func() *stacktrace {
			if hasStack || !opts.withStack {
				return nil
			}

			f := captureStacktrace(3, stacktraceFull)
			return f
		}(),
		causes: causes,
		msg: func() string {
			if skipMsg {
				return ""
			}
			return fmt.Sprintf(msg, args...)
		}(),
		fields: FromCtx(ctx),
	})
}

func applyHook(e *ExtendedError) error {
	if opts.hook == nil {
		return e
	}

	return opts.hook(e)
}

func writeStack(buf *buffer.Buffer, stack []*stacktrace) {
	formatter := newStackFormatter(buf)
	for i, v := range stack {
		if i != 0 {
			buf.WriteString("\n\n" + strings.Repeat("-", 20) + "\n")
		}
		formatter.FormatStack(v)
	}
}

func (e *ExtendedError) frames() (r []*stacktrace) {
	if e.stack != nil {
		return []*stacktrace{e.stack}
	}

	for _, c := range e.causes {
		r = append(r, frames(c)...)
	}

	return r
}

func frames(err error) (r []*stacktrace) {
	switch x := err.(type) {
	case *ExtendedError:
		r = append(r, x.frames()...)
	case interface{ Unwrap() []error }:
		for _, v := range x.Unwrap() {
			r = append(r, frames(v)...)
		}
	case interface{ Unwrap() error }:
		r = append(r, frames(x.Unwrap())...)
	}

	return r
}
