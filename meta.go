package rel

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/lib/pq"
)

var metaCache sync.Map

// ColumnMeta represents a column metadata.
type ColumnMeta struct {
	pk   bool
	name string
	path []int
}

// Identifier returns a quoted column name.
func (cm *ColumnMeta) Identifier() string {
	return pq.QuoteIdentifier(cm.name)
}

// ListColumnMeta is a list of column metadata.
type ListColumnMeta []*ColumnMeta

// Identifiers returns a list of quoted column Names.
func (m ListColumnMeta) Identifiers() []string {
	res := make([]string, len(m))
	for i, cm := range m {
		res[i] = cm.Identifier()
	}
	return res
}

// Names returns a list of unquoted column Names.
func (m ListColumnMeta) Names() []string {
	res := make([]string, len(m))
	for i, cm := range m {
		res[i] = cm.name
	}
	return res
}

// Metadata represents struct metadata.
type Metadata[T any] struct {
	// Primary key strategy (sequence or generated).
	pkStrategy PKStrategy
	// All columns.
	columns ListColumnMeta
	// Primary key columns.
	pkColumns ListColumnMeta
	// columns for insert according to pk strategy.
	insertColumns ListColumnMeta
	// columns for update (all columns except primary key).
	updateColumns ListColumnMeta
	// Map of columns by name for quick access.
	columnsMap map[string]*ColumnMeta
}

func (m Metadata[T]) PkStrategy() PKStrategy {
	return m.pkStrategy
}

func (m Metadata[T]) Columns() ListColumnMeta {
	return m.columns
}

func (m Metadata[T]) PKColumns() ListColumnMeta {
	return m.pkColumns
}

func (m Metadata[T]) InsertColumns() ListColumnMeta {
	return m.insertColumns
}

func (m Metadata[T]) UpdateColumns() ListColumnMeta {
	return m.updateColumns
}

// NewMeta creates a new Metadata instance for the given T type.
func NewMeta[T any](pkStrategy PKStrategy, pk ...string) (*Metadata[T], error) {
	if len(pk) == 0 {
		return nil, errors.New("no primary key specified")
	}

	typeName := reflect.TypeOf((*T)(nil)).Elem().String()
	if cached, ok := metaCache.Load(typeName); ok {
		return cached.(*Metadata[T]), nil
	}

	m := &Metadata[T]{pkStrategy: pkStrategy}
	var err error
	m.columns, err = columnsMeta[T](pk...)
	if err != nil {
		return nil, fmt.Errorf("build columns meta: %w", err)
	}
	m.pkColumns = pkColumns(m.columns)
	if len(m.pkColumns) == 0 {
		return nil, errors.New("no primary key columns found")
	}
	m.insertColumns = insertColumns(m.columns, pkStrategy)
	m.updateColumns = updateColumns(m.columns)
	m.columnsMap = columnsMetaMap(m.columns)

	metaCache.Store(typeName, m)

	return m, nil
}

func columnsMeta[T any](pk ...string) ([]*ColumnMeta, error) {
	var input T
	t := reflect.TypeOf(input)

	// if input is a pointer, get the underlying type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct or pointer to a struct")
	}

	var metas []*ColumnMeta

	// collectMeta recursively collects column metadata.
	var collectMeta func(reflect.Type, []int)
	collectMeta = func(t reflect.Type, parentPath []int) {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldPath := append(parentPath, i)

			// check if the column is anonymous struct
			if field.Anonymous && field.Type.Kind() == reflect.Struct {
				collectMeta(field.Type, fieldPath)
				continue
			}

			if tag, ok := field.Tag.Lookup("db"); ok {
				isPk := false
				for _, p := range pk {
					if p == tag {
						isPk = true
						break
					}
				}
				metas = append(metas, &ColumnMeta{
					pk:   isPk,
					name: tag,
					path: fieldPath,
				})
			}
		}
	}

	collectMeta(t, nil)
	return metas, nil
}

// columnsMetaMap returns a map of column metadata by name.
func columnsMetaMap(columns []*ColumnMeta) map[string]*ColumnMeta {
	cmm := make(map[string]*ColumnMeta)
	for _, cm := range columns {
		cmm[cm.name] = cm
	}
	return cmm
}

// pkColumns returns primary key columns.
func pkColumns(columns []*ColumnMeta) ListColumnMeta {
	return filter(columns, func(cm *ColumnMeta) bool {
		return cm.pk
	})
}

func insertColumns(columns []*ColumnMeta, pkStrategy PKStrategy) ListColumnMeta {
	return filter(columns, func(cm *ColumnMeta) bool {
		return !cm.pk || pkStrategy == PkStrategyGenerated
	})
}

func updateColumns(columns []*ColumnMeta) ListColumnMeta {
	return filter(columns, func(cm *ColumnMeta) bool {
		return !cm.pk
	})
}

func filter(columns ListColumnMeta, f func(*ColumnMeta) bool) ListColumnMeta {
	var res ListColumnMeta
	for _, cm := range columns {
		if f(cm) {
			res = append(res, cm)
		}
	}
	return res
}
