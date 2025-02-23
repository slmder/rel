package db

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type Identifier string

func (i Identifier) Quoted() string {
	return pq.QuoteIdentifier(string(i))
}

type operator int

const (
	opUnknown operator = iota
	opEq
	opNeq
	opIn
	opNotIn
	opAny
	opNotAll
	opIsNull
	opNotNull
	opGt
	opGte
	opLt
	opLte
	opBetween
	opLike
	opLikeLower
)

type Cond []Expr

func (c Cond) Split() ([]any, []string) {
	var args []interface{}
	var expr []string

	for _, e := range c {
		column := pq.QuoteIdentifier(e.field)
		switch e.op {
		case opEq:
			switch {
			case len(e.arg) < 1:
				expr = append(expr, column+" IS NULL")
			default:
				expr = append(expr, column+" = "+ArgsAdd(args, e.arg[0]))
			}
		case opNeq:
			switch {
			case len(e.arg) < 1:
				expr = append(expr, column+" NOT IS NULL")
			default:
				expr = append(expr, column+" <> "+ArgsAdd(args, e.arg[0]))
			}
		case opIn:
			var in []string
			for _, a := range e.arg {
				in = append(in, ArgsAdd(args, a))
			}
			expr = append(expr, column+" IN ("+strings.Join(in, ",")+")")
		case opNotIn:
			var in []string
			for _, a := range e.arg {
				in = append(in, ArgsAdd(args, a))
			}
			expr = append(expr, column+" NOT IN ("+strings.Join(in, ",")+")")
		case opAny:
			expr = append(expr, column+" = ANY "+ArgsAdd(args, pq.Array(e.arg)))
		case opNotAll:
			expr = append(expr, column+" <> ALL "+ArgsAdd(args, pq.Array(e.arg)))
		case opIsNull:
			expr = append(expr, column+" IS NULL")
		case opNotNull:
			expr = append(expr, column+" IS NOT NULL")
		case opGt:
			expr = append(expr, column+" > "+ArgsAdd(args, e.arg[0]))
		case opGte:
			expr = append(expr, column+" >= "+ArgsAdd(args, e.arg[0]))
		case opLt:
			if i, ok := e.arg[0].(Identifier); ok {
				expr = append(expr, column+" < "+i.Quoted())
				continue
			}
			expr = append(expr, column+" < "+ArgsAdd(args, e.arg[0]))
		case opLte:
			if i, ok := e.arg[0].(Identifier); ok {
				expr = append(expr, column+" <= "+i.Quoted())
				continue
			}
			expr = append(expr, column+" <= "+ArgsAdd(args, e.arg[0]))
		case opBetween:
			expr = append(expr, column+" BETWEEN "+ArgsAdd(args, e.arg[0])+" AND "+ArgsAdd(args, e.arg[1]))
		case opLike:
			expr = append(expr, column+" LIKE "+ArgsAdd(args, "%"+fmt.Sprint(e.arg[0])+"%"))
		case opLikeLower:
			expr = append(expr, "LOWER("+column+") LIKE "+ArgsAdd(args, strings.ToLower("%"+fmt.Sprint(e.arg[0])+"%")))
		case opUnknown:
		}
	}

	return args, expr
}

type Expr struct {
	op    operator
	field string
	arg   []any
}

func Eq(field string, arg any) Expr {
	return Expr{
		op:    opEq,
		field: field,
		arg:   []any{arg},
	}
}

func Neq(field string, arg any) Expr {
	return Expr{
		op:    opNeq,
		field: field,
		arg:   []any{arg},
	}
}

func In(field string, args ...any) Expr {
	return Expr{
		op:    opIn,
		field: field,
		arg:   args,
	}
}

func NotIn(field string, args ...any) Expr {
	return Expr{
		op:    opNotIn,
		field: field,
		arg:   args,
	}
}

func Any(field string, args []any) Expr {
	return Expr{
		op:    opAny,
		field: field,
		arg:   args,
	}
}

func NotAny(field string, args []any) Expr {
	return Expr{
		op:    opNotAll,
		field: field,
		arg:   args,
	}
}

func IsNull(field string) Expr {
	return Expr{
		op:    opIsNull,
		field: field,
	}
}

func NotNull(field string, arg any) Expr {
	return Expr{
		op:    opNotNull,
		field: field,
		arg:   []any{arg},
	}
}

func Gt(field string, arg any) Expr {
	return Expr{
		op:    opGt,
		field: field,
		arg:   []any{arg},
	}
}

func Gte(field string, arg any) Expr {
	return Expr{
		op:    opGte,
		field: field,
		arg:   []any{arg},
	}
}

func Lt(field string, arg any) Expr {
	return Expr{
		op:    opLt,
		field: field,
		arg:   []any{arg},
	}
}

func Lte(field string, arg any) Expr {
	return Expr{
		op:    opLte,
		field: field,
		arg:   []any{arg},
	}
}

func Between(field string, a, b any) Expr {
	return Expr{
		op:    opBetween,
		field: field,
		arg:   []any{a, b},
	}
}

func Like(field string, a string) Expr {
	return Expr{
		op:    opLike,
		field: field,
		arg:   []any{a},
	}
}

func LikeLower(field string, a string) Expr {
	return Expr{
		op:    opLikeLower,
		field: field,
		arg:   []any{a},
	}
}
