package qbuilder

import (
	"strings"
)

const insertBufferInitialGrowBytes = 200

type InsertBuilder struct {
	tableName          string
	alias              string
	columns            []string
	values             [][]string
	args               []interface{}
	conflictTarget     string
	conflictConstraint bool
	conflictAction     string
	conflictSet        map[string]string
	returning          []string
}

func Insert(into string) *InsertBuilder {
	builder := InsertBuilder{}
	return builder.Insert(into)
}

func (b *InsertBuilder) Insert(into string) *InsertBuilder {
	b.tableName = into
	return b
}

func (b *InsertBuilder) As(alias string) *InsertBuilder {
	b.alias = alias
	return b
}

func (b *InsertBuilder) Columns(cols ...string) *InsertBuilder {
	b.columns = cols
	return b
}

func (b *InsertBuilder) Values(values ...[]string) *InsertBuilder {
	b.values = values
	return b
}

func (b *InsertBuilder) OnConflict(target string, constraint bool) *InsertBuilder {
	b.conflictTarget = target
	b.conflictConstraint = constraint
	return b
}

func (b *InsertBuilder) DoNothing() *InsertBuilder {
	b.conflictAction = "DO NOTHING"
	return b
}

func (b *InsertBuilder) DoUpdate(kv map[string]string) *InsertBuilder {
	b.conflictAction = "DO UPDATE"
	b.conflictSet = kv
	return b
}

func (b *InsertBuilder) Returning(alias ...string) *InsertBuilder {
	b.returning = alias
	return b
}

func (b *InsertBuilder) ToSQL() string {
	var out strings.Builder
	out.Grow(insertBufferInitialGrowBytes)
	out.WriteString("INSERT INTO ")
	out.WriteString(b.tableName)
	comma := ", "
	if b.alias != "" {
		out.WriteString(" AS ")
		out.WriteString(b.alias)
	}
	if len(b.columns) > 0 {
		out.WriteString(" (")
		out.WriteString(b.columns[0])
		for _, s := range b.columns[1:] {
			out.WriteString(comma)
			out.WriteString(s)
		}
		out.WriteString(")")
	}
	if len(b.values) > 0 {
		out.WriteString(" VALUES ")
		out.WriteString("(")
		out.WriteString(strings.Join(b.values[0], ", "))
		out.WriteString(")")
		for _, s := range b.values[1:] {
			out.WriteString(comma)
			out.WriteString("(")
			out.WriteString(strings.Join(s, ", "))
			out.WriteString(")")
		}
	}

	if b.conflictTarget != "" {
		out.WriteString(" ON CONFLICT ")
		if b.conflictConstraint {
			out.WriteString("ON CONSTRAINT ")
		}
		out.WriteString("(")
		out.WriteString(b.conflictTarget)
		out.WriteString(") ")
		out.WriteString(b.conflictAction)
		if len(b.conflictSet) > 0 {
			out.WriteString(" SET ")
			i := len(b.conflictSet) - 1
			for c, v := range b.conflictSet {
				out.WriteString(c)
				out.WriteString(" = ")
				out.WriteString(v)
				if i != 0 {
					out.WriteString(comma)
				}
				i--
			}
		}
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
