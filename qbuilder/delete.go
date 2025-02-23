package qbuilder

import (
	"fmt"
	"strings"
)

const deleteBufferInitialGrowBytes = 200

type DeleteBuilder struct {
	tableName string
	alias     string
	usingExpr []string
	whereExpr []string
	returning []string
}

func Delete(rel string) *DeleteBuilder {
	builder := DeleteBuilder{}
	return builder.Delete(rel)
}

func (b *DeleteBuilder) Delete(from string) *DeleteBuilder {
	b.tableName = from
	return b
}

func (b *DeleteBuilder) As(alias string) *DeleteBuilder {
	b.alias = alias
	return b
}

func (b *DeleteBuilder) Using(expr ...string) *DeleteBuilder {
	b.usingExpr = expr
	return b
}

func (b *DeleteBuilder) AndUsing(expr string, a ...interface{}) *DeleteBuilder {
	if len(a) > 0 {
		b.usingExpr = append(b.usingExpr, fmt.Sprintf(expr, a...))
	} else {
		b.usingExpr = append(b.usingExpr, expr)
	}
	return b
}

func (b *DeleteBuilder) ResetUsing() *DeleteBuilder {
	b.usingExpr = []string{}
	return b
}

func (b *DeleteBuilder) Where(expr string, a ...interface{}) *DeleteBuilder {
	if len(b.whereExpr) > 0 {
		b.whereExpr = []string{}
	}
	return b.AndWhere(expr, a...)
}

func (b *DeleteBuilder) AndWhere(expr string, a ...interface{}) *DeleteBuilder {
	if len(a) > 0 {
		b.whereExpr = append(b.whereExpr, fmt.Sprintf(expr, a...))
	} else {
		b.whereExpr = append(b.whereExpr, expr)
	}
	return b
}

func (b *DeleteBuilder) Returning(alias ...string) *DeleteBuilder {
	b.returning = alias
	return b
}

func (b *DeleteBuilder) String() string {
	return b.ToSQL()
}

func (b *DeleteBuilder) ToSQL() string {
	var out strings.Builder
	out.Grow(deleteBufferInitialGrowBytes)
	out.WriteString("DELETE FROM ")
	out.WriteString(b.tableName)
	comma := ", "
	if b.alias != "" {
		out.WriteString(" AS ")
		out.WriteString(b.alias)
	}
	if len(b.usingExpr) > 0 {
		out.WriteString(" USING ")
		out.WriteString(b.usingExpr[0])
		for _, s := range b.usingExpr[1:] {
			out.WriteString(comma)
			out.WriteString(s)
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
