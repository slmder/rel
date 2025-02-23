package db

import (
	"fmt"
	"reflect"
)

type scanFunc func(dest ...any) error

// scanRow scans all struct fields using dst fields as a destination.
func scanRow[T any](scan scanFunc, m *Metadata[T], dst *T) error {
	pointers, err := getFieldsPointers[T](m.Columns.Names(), m, dst)
	if err != nil {
		return err
	}
	return scan(pointers...)
}

// scanPK scans only primary key fields using dst fields as a destination.
func scanPK[T any](scan scanFunc, m *Metadata[T], dst *T) error {
	dest, err := getFieldsPointers[T](m.PKColumns.Names(), m, dst)
	if err != nil {
		return err
	}
	return scan(dest...)
}

// getFieldsPointers returns pointers to fields of the struct.
func getFieldsPointers[T any](fields []string, meta *Metadata[T], entity *T) ([]any, error) {
	v := reflect.ValueOf(entity)

	// If a pointer is passed, we get the value it points to.
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity must be a struct or pointer to a struct")
	}

	result := make([]any, len(fields))
	// Retrieving field pointers.
	for i, name := range fields {
		if fieldInfo, found := meta.columnsMap[name]; found {
			fieldValue := v
			for _, index := range fieldInfo.path {
				fieldValue = fieldValue.Field(index)
			}
			result[i] = fieldValue.Addr().Interface()
		} else {
			return nil, fmt.Errorf("field not found: %s", name)
		}
	}

	return result, nil
}
