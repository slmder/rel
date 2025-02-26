package rel

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
		column := pq.QuoteIdentifier(e.column)
		switch e.op {
		case opEq:
			switch {
			case len(e.arg) < 1:
				expr = append(expr, column+" IS NULL")
			default:
				if i, ok := e.arg[0].(Identifier); ok {
					expr = append(expr, column+" = "+i.Quoted())
					continue
				}
				expr = append(expr, column+" = "+ArgsAdd(&args, e.arg[0]))
			}
		case opNeq:
			switch {
			case len(e.arg) < 1:
				expr = append(expr, column+" IS NOT NULL")
			default:
				if i, ok := e.arg[0].(Identifier); ok {
					expr = append(expr, column+" <> "+i.Quoted())
					continue
				}
				expr = append(expr, column+" <> "+ArgsAdd(&args, e.arg[0]))
			}
		case opIn:
			var in []string
			for _, a := range e.arg {
				in = append(in, ArgsAdd(&args, a))
			}
			expr = append(expr, column+" IN ("+strings.Join(in, ",")+")")
		case opNotIn:
			var in []string
			for _, a := range e.arg {
				in = append(in, ArgsAdd(&args, a))
			}
			expr = append(expr, column+" NOT IN ("+strings.Join(in, ",")+")")
		case opAny:
			if i, ok := e.arg[0].(Identifier); ok {
				expr = append(expr, column+" = ANY "+i.Quoted())
				continue
			}
			expr = append(expr, column+" = ANY "+ArgsAdd(&args, pq.Array(e.arg)))
		case opNotAll:
			if i, ok := e.arg[0].(Identifier); ok {
				expr = append(expr, column+" <> ALL "+i.Quoted())
				continue
			}
			expr = append(expr, column+" <> ALL "+ArgsAdd(&args, pq.Array(e.arg)))
		case opIsNull:
			expr = append(expr, column+" IS NULL")
		case opNotNull:
			expr = append(expr, column+" IS NOT NULL")
		case opGt:
			if i, ok := e.arg[0].(Identifier); ok {
				expr = append(expr, column+" > "+i.Quoted())
				continue
			}
			expr = append(expr, column+" > "+ArgsAdd(&args, e.arg[0]))
		case opGte:
			if i, ok := e.arg[0].(Identifier); ok {
				expr = append(expr, column+" >= "+i.Quoted())
				continue
			}
			expr = append(expr, column+" >= "+ArgsAdd(&args, e.arg[0]))
		case opLt:
			if i, ok := e.arg[0].(Identifier); ok {
				expr = append(expr, column+" < "+i.Quoted())
				continue
			}
			expr = append(expr, column+" < "+ArgsAdd(&args, e.arg[0]))
		case opLte:
			if i, ok := e.arg[0].(Identifier); ok {
				expr = append(expr, column+" <= "+i.Quoted())
				continue
			}
			expr = append(expr, column+" <= "+ArgsAdd(&args, e.arg[0]))
		case opBetween:
			a, aok := e.arg[0].(Identifier)
			b, bok := e.arg[1].(Identifier)
			if aok && bok {
				expr = append(expr, column+" BETWEEN "+a.Quoted()+" AND "+b.Quoted())
				continue
			}
			expr = append(expr, column+" BETWEEN "+ArgsAdd(&args, e.arg[0])+" AND "+ArgsAdd(&args, e.arg[1]))
		case opLike:
			if len(e.arg) < 1 || e.arg[0] == nil {
				expr = append(expr, column+" = ''")
				continue
			}
			expr = append(expr, column+" LIKE "+ArgsAdd(&args, "%"+fmt.Sprint(e.arg[0])+"%"))
		case opLikeLower:
			if len(e.arg) < 1 || e.arg[0] == nil {
				expr = append(expr, column+" = ''")
				continue
			}
			expr = append(expr, "LOWER("+column+") LIKE "+ArgsAdd(&args, strings.ToLower("%"+fmt.Sprint(e.arg[0])+"%")))
		case opUnknown:
		}
	}

	return args, expr
}

// ArgsAdd adds an argument to the slice and returns a placeholder for it.
func ArgsAdd(args *[]any, arg any) string {
	*args = append(*args, arg)
	return fmt.Sprintf("$%d", len(*args))
}

type Expr struct {
	op     operator
	column string
	arg    []any
}

func Eq(column string, arg any) Expr {
	return Expr{
		op:     opEq,
		column: column,
		arg:    []any{arg},
	}
}

func Neq(column string, arg any) Expr {
	return Expr{
		op:     opNeq,
		column: column,
		arg:    []any{arg},
	}
}

func In(column string, args ...any) Expr {
	return Expr{
		op:     opIn,
		column: column,
		arg:    args,
	}
}

func NotIn(column string, args ...any) Expr {
	return Expr{
		op:     opNotIn,
		column: column,
		arg:    args,
	}
}

func Any(column string, args ...any) Expr {
	return Expr{
		op:     opAny,
		column: column,
		arg:    args,
	}
}

func NotAll(column string, args ...any) Expr {
	return Expr{
		op:     opNotAll,
		column: column,
		arg:    args,
	}
}

func IsNull(column string) Expr {
	return Expr{
		op:     opIsNull,
		column: column,
	}
}

func NotNull(column string) Expr {
	return Expr{
		op:     opNotNull,
		column: column,
	}
}

func Gt(column string, arg any) Expr {
	return Expr{
		op:     opGt,
		column: column,
		arg:    []any{arg},
	}
}

func Gte(column string, arg any) Expr {
	return Expr{
		op:     opGte,
		column: column,
		arg:    []any{arg},
	}
}

func Lt(column string, arg any) Expr {
	return Expr{
		op:     opLt,
		column: column,
		arg:    []any{arg},
	}
}

func Lte(column string, arg any) Expr {
	return Expr{
		op:     opLte,
		column: column,
		arg:    []any{arg},
	}
}

func Between(column string, a, b any) Expr {
	return Expr{
		op:     opBetween,
		column: column,
		arg:    []any{a, b},
	}
}

func Like(column string, a string) Expr {
	return Expr{
		op:     opLike,
		column: column,
		arg:    []any{a},
	}
}

func LikeLower(column string, a string) Expr {
	return Expr{
		op:     opLikeLower,
		column: column,
		arg:    []any{a},
	}
}
