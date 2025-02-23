package qbuilder

import (
	"fmt"
	"strings"
)

const updateBufferInitialGrowBytes = 200

type UpdateBuilder struct {
	tableName string
	alias     string
	set       map[string]string
	whereExpr []string
	returning []string
}

func Update(rel string) *UpdateBuilder {
	builder := UpdateBuilder{}
	return builder.Update(rel)
}

func (b *UpdateBuilder) Update(table string) *UpdateBuilder {
	b.tableName = table
	return b
}

func (b *UpdateBuilder) As(alias string) *UpdateBuilder {
	b.alias = alias
	return b
}

func (b *UpdateBuilder) Set(col, val string) *UpdateBuilder {
	if b.set == nil {
		b.set = make(map[string]string)
	}
	b.set[col] = val
	return b
}

func (b *UpdateBuilder) SetMap(m map[string]string) *UpdateBuilder {
	b.set = m
	return b
}

func (b *UpdateBuilder) Where(expr string, a ...interface{}) *UpdateBuilder {
	if len(b.whereExpr) > 0 {
		b.whereExpr = []string{}
	}
	return b.AndWhere(expr, a...)
}

func (b *UpdateBuilder) AndWhere(expr string, a ...interface{}) *UpdateBuilder {
	if len(a) > 0 {
		b.whereExpr = append(b.whereExpr, fmt.Sprintf(expr, a...))
	} else {
		b.whereExpr = append(b.whereExpr, expr)
	}
	return b
}

func (b *UpdateBuilder) Returning(columns ...string) *UpdateBuilder {
	b.returning = columns
	return b
}

func (b *UpdateBuilder) String() string {
	return b.ToSQL()
}

func (b *UpdateBuilder) ToSQL() string { // nolint:funlen
	var out strings.Builder
	out.Grow(updateBufferInitialGrowBytes)
	out.WriteString("UPDATE ")
	out.WriteString(b.tableName)
	comma := ", "
	if b.alias != "" {
		out.WriteString(" AS ")
		out.WriteString(b.alias)
	}
	if len(b.set) > 0 {
		out.WriteString(" SET ")
		i := len(b.set) - 1
		for c, v := range b.set {
			out.WriteString(c)
			out.WriteString(" = ")
			out.WriteString(v)
			if i != 0 {
				out.WriteString(comma)
			}
			i--
		}
	}
	if len(b.whereExpr) > 0 {
		and := " AND "
		out.WriteString(" WHERE ")
		out.WriteString("(")
		out.WriteString("(")
		out.WriteString(b.whereExpr[0])
		out.WriteString(")")
		for _, s := range b.whereExpr[1:] {
			out.WriteString(and)
			out.WriteString("(")
			out.WriteString(s)
			out.WriteString(")")
		}
		out.WriteString(")")
	}

	if len(b.returning) > 0 {
		out.WriteString(" RETURNING ")
		out.WriteString(b.returning[0])
		for _, s := range b.returning[1:] {
			out.WriteString(comma)
			out.WriteString(s)
		}
	}
	return out.String()
}
