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

func (b *CallBuilder) Argument(expr string, a ...interface{}) *CallBuilder {
	return b.AddArgument(expr, a...)
}

func (b *CallBuilder) Arguments(args ...string) *CallBuilder {
	for _, arg := range args {
		b.arguments = append(b.arguments, arg)
	}
	return b
}

func (b *CallBuilder) AddArgument(expr string, a ...interface{}) *CallBuilder {
	if len(a) > 0 {
		b.arguments = append(b.arguments, fmt.Sprintf(expr, a...))
	} else {
		b.arguments = append(b.arguments, expr)
	}
	return b
}

func (b *CallBuilder) Clear() *CallBuilder {
	b.arguments = nil
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
		for i, arg := range b.arguments {
			if i > 0 {
				out.WriteString(", ")
			}
			out.WriteString(fmt.Sprintf("%v", arg))
		}
		out.WriteString(")")
	}
	return out.String()
}
