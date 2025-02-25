package qbuilder

import (
	"fmt"
	"strings"
)

const updateBufferInitialGrowBytes = 200

type UpdateBuilder struct {
	tableName  string
	alias      string
	setColumns []string
	setValues  []string
	whereExpr  []string
	returning  []string
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

func (b *UpdateBuilder) Set(column, val string) *UpdateBuilder {
	b.setColumns = append(b.setColumns, column)
	b.setValues = append(b.setValues, val)
	return b
}

func (b *UpdateBuilder) SetMap(m map[string]string) *UpdateBuilder {
	for col, val := range m {
		b.Set(col, val)
	}
	return b
}

func (b *UpdateBuilder) Where(expr string, a ...interface{}) *UpdateBuilder {
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
	if len(b.setColumns) > 0 {
		out.WriteString(" SET ")
		i := len(b.setColumns) - 1
		for j, c := range b.setColumns {
			out.WriteString(c)
			out.WriteString(" = ")
			out.WriteString(b.setValues[j])
			if i != 0 {
				out.WriteString(comma)
			}
			i--
		}
	}
	if len(b.whereExpr) > 0 {
		out.WriteString(" WHERE " + strings.Join(b.whereExpr, " AND "))
	}

	if len(b.returning) > 0 {
		out.WriteString(" RETURNING " + strings.Join(b.returning, ", "))
	}

	return out.String()
}
