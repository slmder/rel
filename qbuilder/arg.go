package qbuilder

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/lib/pq"
)

var (
	QuoteLiteral    = pq.QuoteLiteral
	QuoteIdentifier = pq.QuoteIdentifier
	NoopQuote       = noopQuoter
)

func noopQuoter(l string) string {
	return l
}

func In[T any](a []T, quote func(l string) string) string {
	if len(a) == 0 {
		return ""
	}
	if quote == nil {
		quote = NoopQuote
	}
	var b strings.Builder
	b.WriteString(quote(fmt.Sprintf("%v", a[0])))
	for _, s := range a {
		b.WriteString(",")
		b.WriteString(quote(fmt.Sprintf("%v", s)))
	}
	return b.String()
}

func InUUID(a []string) string {
	return In(a, QuoteLiteral)
}

func InString(a []string) string {
	return In(a, QuoteLiteral)
}

func Like(a interface{}) string {
	return QuoteLiteral("%" + fmt.Sprint(a) + "%")
}

func LikeLower(a string) string {
	return Like(strings.ToLower(a))
}

func String(a string) string {
	return QuoteLiteral(a)
}

func Time(a time.Time) string {
	return QuoteLiteral(string(pq.FormatTimestamp(a)))
}

func StructTaggedFields(t reflect.Type, names *[]string, dbnames *[]string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag, ok := field.Tag.Lookup("db"); ok {
			if tag != "" {
				*names = append(*names, field.Name)
				*dbnames = append(*dbnames, tag)
			}
		} else if field.Type.Kind() == reflect.Struct && field.Anonymous {
			StructTaggedFields(field.Type, names, dbnames)
		}
	}
}
