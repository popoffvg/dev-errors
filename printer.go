package errors

import (
	"fmt"

	"github.com/popoffvg/dev-errors/internal/bufferpool"
)

type (
	printer interface {
		Print(msg string, stack []*stacktrace, fs []Field) string
	}

	defaultPrinter struct{}

	verbosePrinter struct{}
)

func (p *defaultPrinter) Print(msg string, stack []*stacktrace, fs []Field) string {
	return msg
}

func (p *verbosePrinter) Print(msg string, stack []*stacktrace, fs []Field) string {
	var buf = bufferpool.Get()

	buf.WriteString("msg:")
	buf.WriteString(msg)
	buf.WriteString("\n")

	if len(fs) > 0 {
		buf.WriteString("fields:")
	}
	for _, f := range fs {
		buf.WriteString("\t")
		buf.WriteString(f.Key)
		buf.WriteString(":")
		fmt.Fprintf(buf, "%+v\n", f.Value)
	}

	if opts.withStack {
		if stack != nil {
			buf.WriteString("stack:\n")
			writeStack(buf, stack)
			buf.WriteString("\n")
		}
	}

	return buf.String()
}
