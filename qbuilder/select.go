package qbuilder

import (
	"fmt"
	"strconv"
	"strings"
)

const selectBufferInitialGrowBytes = 200

type SelectBuilder struct {
	distinct     bool
	subSelect    bool
	selectExpr   []string
	alias        string
	fromExpr     string
	join         []string
	whereExpr    []string
	groupByExpr  []string
	havingExpr   []string
	unionExpr    []string
	unionAllExpr []string
	orderByExpr  []string
	limit        string
	offset       string
	forExpr      string
}

func Select(columns ...string) *SelectBuilder {
	builder := SelectBuilder{}
	return builder.Select(columns...)
}

func SubSelect(columns ...string) *SelectBuilder {
	builder := SelectBuilder{
		subSelect: true,
	}
	return builder.Select(columns...)
}

func (b *SelectBuilder) Select(expr ...string) *SelectBuilder {
	b.selectExpr = expr
	return b
}

func (b *SelectBuilder) AddSelect(expr ...string) *SelectBuilder {
	b.selectExpr = append(b.selectExpr, expr...)
	return b
}

func (b *SelectBuilder) Distinct(v bool) *SelectBuilder {
	b.distinct = v
	return b
}

func (b *SelectBuilder) From(rel string, alias ...string) *SelectBuilder {
	b.fromExpr = rel
	if len(alias) > 0 {
		b.As(alias[0])
	}
	return b
}

func (b *SelectBuilder) As(alias string) *SelectBuilder {
	b.alias = alias
	return b
}

func (b *SelectBuilder) LeftJoin(rel, alias, cond string) *SelectBuilder {
	return b.Join(JoinTypeLeft, rel, alias, cond)
}

func (b *SelectBuilder) RightJoin(rel, alias, cond string) *SelectBuilder {
	return b.Join(joinTypeRight, rel, alias, cond)
}

func (b *SelectBuilder) CrossJoin(rel, alias, cond string) *SelectBuilder {
	return b.Join(joinTypeCross, rel, alias, cond)
}

func (b *SelectBuilder) InnerJoin(rel, alias, cond string) *SelectBuilder {
	return b.Join(joinTypeInner, rel, alias, cond)
}

func (b *SelectBuilder) Join(dir joinType, rel, alias, cond string) *SelectBuilder {
	b.join = append(b.join, dir.String()+" JOIN "+rel+" AS "+alias+" ON "+cond)
	return b
}

func (b *SelectBuilder) ResetJoin() *SelectBuilder {
	b.join = []string{}
	return b
}

func (b *SelectBuilder) Where(expr string, a ...interface{}) *SelectBuilder {
	if len(b.whereExpr) > 0 {
		b.whereExpr = []string{}
	}
	return b.AndWhere(expr, a...)
}

func (b *SelectBuilder) AndWhere(expr string, a ...interface{}) *SelectBuilder {
	if len(a) > 0 {
		b.whereExpr = append(b.whereExpr, fmt.Sprintf(expr, a...))
	} else {
		b.whereExpr = append(b.whereExpr, expr)
	}
	return b
}

func (b *SelectBuilder) GroupBy(expr ...string) *SelectBuilder {
	b.groupByExpr = expr
	return b
}

func (b *SelectBuilder) AndGroupBy(expr string) *SelectBuilder {
	b.groupByExpr = append(b.groupByExpr, expr)
	return b
}

func (b *SelectBuilder) ResetGroupBy() *SelectBuilder {
	b.groupByExpr = []string{}
	return b
}

func (b *SelectBuilder) Having(expr ...string) *SelectBuilder {
	b.havingExpr = expr
	return b
}

func (b *SelectBuilder) AndHaving(expr string) *SelectBuilder {
	b.havingExpr = append(b.havingExpr, expr)
	return b
}

func (b *SelectBuilder) Union(expr string) *SelectBuilder {
	b.unionExpr = []string{expr}
	return b
}

func (b *SelectBuilder) AndUnion(expr string) *SelectBuilder {
	b.unionExpr = append(b.unionExpr, expr)
	return b
}

func (b *SelectBuilder) UnionAll(expr string) *SelectBuilder {
	b.unionAllExpr = []string{expr}
	return b
}

func (b *SelectBuilder) AndUnionAll(expr string) *SelectBuilder {
	b.unionAllExpr = append(b.unionAllExpr, expr)
	return b
}

func (b *SelectBuilder) ResetHaving() *SelectBuilder {
	b.havingExpr = []string{}
	return b
}

func (b *SelectBuilder) OrderBy(col string, dir SortDirection) *SelectBuilder {
	b.orderByExpr = []string{col + " " + dir.String()}
	return b
}

func (b *SelectBuilder) AndOrderBy(col string, dir SortDirection) *SelectBuilder {
	b.orderByExpr = append(b.orderByExpr, col+" "+dir.String())
	return b
}

func (b *SelectBuilder) ResetOrderBy() *SelectBuilder {
	b.orderByExpr = []string{}
	return b
}

func (b *SelectBuilder) Limit(limit uint32) *SelectBuilder {
	if limit > 0 {
		b.limit = strconv.Itoa(int(limit))
	} else {
		b.limit = ""
	}
	return b
}

func (b *SelectBuilder) Offset(offset uint32) *SelectBuilder {
	if offset > 0 {
		b.offset = strconv.Itoa(int(offset))
	} else {
		b.offset = ""
	}
	return b
}

func (b *SelectBuilder) For(mode rowLevelLockMode) *SelectBuilder {
	b.forExpr = mode.String()
	return b
}

func (b *SelectBuilder) String() string {
	return b.ToSQL()
}

func (b *SelectBuilder) ToSQL() string { // nolint:funlen
	var out strings.Builder
	out.Grow(selectBufferInitialGrowBytes)
	if b.subSelect {
		out.WriteString("(")
	}
	out.WriteString("SELECT ")
	if b.distinct {
		out.WriteString("DISTINCT ")
	}
	comma := ", "
	switch len(b.selectExpr) {
	case 0:
		out.WriteString("*")
	case 1:
		out.WriteString(b.selectExpr[0])
	default:
		out.WriteString(b.selectExpr[0])
		for _, s := range b.selectExpr[1:] {
			out.WriteString(comma)
			out.WriteString(s)
		}
	}
	out.WriteString(" FROM ")
	out.WriteString(b.fromExpr)

	if b.alias != "" {
		out.WriteString(" AS ")
		out.WriteString(b.alias)
		out.WriteString(" ")
	}

	if len(b.join) > 0 {
		for _, s := range b.join {
			out.WriteString(s)
			out.WriteString(" ")
		}
	}

	if len(b.whereExpr) > 0 {
		and := " AND "
		out.WriteString(" WHERE (")
		out.WriteString(b.whereExpr[0])
		out.WriteString(")")
		for _, s := range b.whereExpr[1:] {
			out.WriteString(and)
			out.WriteString("(")
			out.WriteString(s)
			out.WriteString(")")
		}
	}

	if len(b.groupByExpr) > 0 {
		out.WriteString(" GROUP BY ")
		out.WriteString(b.groupByExpr[0])
		for _, s := range b.groupByExpr[1:] {
			out.WriteString(comma)
			out.WriteString(s)
		}
	}

	if len(b.havingExpr) > 0 {
		out.WriteString(" HAVING ")
		out.WriteString(b.havingExpr[0])
		for _, s := range b.havingExpr[1:] {
			out.WriteString(comma)
			out.WriteString(s)
		}
	}

	if len(b.unionExpr) > 0 {
		for _, s := range b.havingExpr {
			out.WriteString(" UNION (")
			out.WriteString(s)
			out.WriteString(")")
		}
	}

	if len(b.unionAllExpr) > 0 {
		for _, s := range b.unionAllExpr {
			out.WriteString(" UNION ALL (")
			out.WriteString(s)
			out.WriteString(")")
		}
	}

	if len(b.orderByExpr) > 0 {
		out.WriteString(" ORDER BY ")
		out.WriteString(b.orderByExpr[0])
		for _, s := range b.orderByExpr[1:] {
			out.WriteString(comma)
			out.WriteString(s)
		}
	}

	if b.limit != "" {
		out.WriteString(" LIMIT ")
		out.WriteString(b.limit)
	}

	if b.offset != "" {
		out.WriteString(" OFFSET ")
		out.WriteString(b.offset)
	}

	if b.forExpr != "" {
		out.WriteString(" FOR ")
		out.WriteString(b.forExpr)
	}
	if b.subSelect {
		out.WriteString(")")
	}
	return out.String()
}

func AndX(exp ...string) string {
	return strings.Join(exp, " AND ")
}
