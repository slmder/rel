package rel

import (
	"fmt"
	"reflect"
)

type scanFunc func(dest ...any) error

// scanRow scans all struct fields using dst fields as a destination.
func scanRow[T any](scan scanFunc, m *Metadata[T], dst *T) error {
	pointers, err := getFieldsPointers[T](m.Columns().Names(), m, dst)
	if err != nil {
		return err
	}
	return scan(pointers...)
}

// getFieldsPointers returns pointers to fields of the struct.
func getFieldsPointers[T any](fields []string, meta *Metadata[T], entity *T) ([]any, error) {
	v := reflect.ValueOf(entity)

	// If a pointer is passed, we get the value it points to.
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, fmt.Errorf("nil entity")
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity must be a struct or pointer to a struct")
	}

	if meta.columnsMap == nil {
		return nil, fmt.Errorf("metadata not initialized")
	}

	result := make([]any, len(fields))
	// Retrieving column pointers.
	for i, name := range fields {
		if fieldInfo, found := meta.columnsMap[name]; found {
			fieldValue := v
			for _, index := range fieldInfo.path {
				fieldValue = fieldValue.Field(index)
			}
			result[i] = fieldValue.Addr().Interface()
		} else {
			return nil, fmt.Errorf("column not found: %s", name)
		}
	}

	return result, nil
}
