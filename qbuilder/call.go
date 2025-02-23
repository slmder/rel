package qbuilder

import (
	"fmt"
	"strings"
)

type CallBuilder struct {
	procName  string
	arguments []string
}

func (b *CallBuilder) Call(proc string) *CallBuilder {
	b.procName = proc
	return b
}

func (b *CallBuilder) Argument(arg string, a ...interface{}) *CallBuilder {
	if len(a) > 0 {
		b.arguments = []string{fmt.Sprintf(arg, a...)}
	} else {
		b.arguments = []string{arg}
	}
	return b
}

func (b *CallBuilder) Arguments(arg ...string) *CallBuilder {
	b.arguments = arg
	return b
}

func (b *CallBuilder) AndArgument(expr string, a ...interface{}) *CallBuilder {
	if len(a) > 0 {
		b.arguments = append(b.arguments, fmt.Sprintf(expr, a...))
	} else {
		b.arguments = append(b.arguments, expr)
	}
	return b
}

func (b *CallBuilder) String() string {
	return b.ToSQL()
}

func (b *CallBuilder) ToSQL() string {
	var out strings.Builder
	out.WriteString("CALL ")
	out.WriteString(b.procName)
	if len(b.arguments) > 0 {
		out.WriteString(" (")
		out.WriteString(b.arguments[0])
		for _, s := range b.arguments[1:] {
			out.WriteString(", ")
			out.WriteString(s)
		}
		out.WriteString(")")
	}

	return out.String()
}
