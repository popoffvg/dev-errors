package errors

import (
	"context"
	"errors"
	"fmt"

	"github.com/popoffvg/dev-errors/fields"
	"github.com/popoffvg/dev-errors/internal/bufferpool"
)

type (
	ExtendedError struct {
		f *stacktrace

		causes []error
		msg    string
		fields []fields.Field
	}
)

// New create new error with stack according to stack strategy
// without field from context.
func New(msg string, args ...any) error {
	return NewCtx(context.Background(), msg, args...)
}

// New create new error with stack according to stack strategy
// with field from context.
//
// If error with frame exists in args then frame will not be added.
func NewCtx(ctx context.Context, msg string, args ...any) error {
	var (
		causes []error

		hasStack       bool
		causeWithStack = new(ExtendedError)
	)

	if withStack {
		for _, v := range args {
			if e, ok := v.(error); ok {
				causes = append(causes, e)
				if !hasStack && errors.As(e, &causeWithStack) && causeWithStack.f != nil {
					hasStack = true
				}
			}
		}
	}

	return &ExtendedError{
		f: func() *stacktrace {
			if hasStack || !withStack {
				return nil
			}

			f := captureStacktrace(2, stacktraceFull)
			return f
		}(),
		causes: causes,
		msg:    fmt.Sprintf(msg, args...),
		fields: fields.FromCtx(ctx),
	}
}

func (e *ExtendedError) Unwrap() error {
	if len(e.causes) == 0 {
		return nil
	}

	return e.causes[len(e.causes)-1]
}

func (e *ExtendedError) Error() string {
	var buf = bufferpool.Get()
	// TODO:: (popoffvg) custom printer

	buf.WriteString("msg:")
	buf.WriteString(e.msg)
	buf.WriteString("\n")

	if len(e.fields) > 0 {
		buf.WriteString("fields:")
	}
	for _, f := range e.fields {
		buf.WriteString("\t")
		buf.WriteString(f.Key)
		buf.WriteString(":")
		// TODO:: (popoffvg) optimize
		buf.WriteString(fmt.Sprintf("%+v\n", f.Value))
	}

	if withStack {
		f := e.frame()
		if f != nil {
			buf.WriteString("stack:\n")
			formatter := newStackFormatter(buf)
			formatter.FormatStack(f)
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

func (e *ExtendedError) Fields() []fields.Field {
	return e.fields
}

func (e *ExtendedError) frame() *stacktrace {
	if e.f != nil {
		return e.f
	}

	var subErr *ExtendedError
	for _, v := range e.causes {
		if errors.As(v, &subErr) {
			return subErr.frame()
		}
	}

	return nil
}
