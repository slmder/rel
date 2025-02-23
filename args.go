package db

import (
	"fmt"
	"reflect"
)

// ArgsAdd adds an argument to the slice and returns a placeholder for it.
func ArgsAdd(args []any, arg any) string {
	args = append(args, arg)
	return fmt.Sprintf("$%d", len(args))
}

// getFieldsValues returns values of the fields of the struct.
func getFieldsValues[T any](fields []string, meta *Metadata[T], entity *T) []any {
	v := reflect.ValueOf(entity)

	// If a pointer is passed, we get the value it points to.
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("entity must be a struct or pointer to a struct")
	}

	result := make([]any, len(fields))
	// Retrieving field values.
	for i, name := range fields {
		if fieldInfo, found := meta.columnsMap[name]; found {
			fieldValue := v
			for _, index := range fieldInfo.path {
				fieldValue = fieldValue.Field(index)
			}
			result[i] = fieldValue.Interface()
		} else {
			panic(fmt.Sprintf("field not found: %s", name))
		}

	}
	return result
}

// getArgsPlaceholders returns given number of placeholders for SQL query.
func getArgsPlaceholders(n int) []string {
	res := make([]string, n)
	for i := range res {
		res[i] = fmt.Sprintf("$%d", i+1)
	}
	return res
}
