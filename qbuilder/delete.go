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

func Delete(from string) *DeleteBuilder {
	builder := DeleteBuilder{}
	return builder.Delete(from)
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
	if b.alias != "" {
		out.WriteString(" AS ")
		out.WriteString(b.alias)
	}
	if len(b.usingExpr) > 0 {
		out.WriteString(" USING " + strings.Join(b.usingExpr, ", "))
	}
	if len(b.whereExpr) > 0 {
		out.WriteString(" WHERE " + strings.Join(b.whereExpr, " AND "))
	}
	if len(b.returning) > 0 {
		out.WriteString(" RETURNING " + strings.Join(b.returning, ", "))
	}

	return out.String()
}
