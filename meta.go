package rel

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/lib/pq"
)

const maxRecursionDepth = 10

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
	if len(m) == 0 {
		return nil
	}
	res := make([]string, len(m))
	for i, cm := range m {
		if cm == nil {
			res[i] = ""
			continue
		}
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

	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Struct && !(t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct) {
		return nil, fmt.Errorf("type %s is not a struct or a pointer to a struct", t.String())
	}

	typeName := t.String()
	if cached, ok := metaCache.Load(typeName); ok {
		if metadata, ok := cached.(*Metadata[T]); ok {
			return metadata, nil
		}
		metaCache.Delete(typeName)
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

	actual, loaded := metaCache.LoadOrStore(typeName, m)
	if loaded {
		return actual.(*Metadata[T]), nil
	}

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
	var collectMeta func(reflect.Type, []int, int) error
	collectMeta = func(t reflect.Type, parentPath []int, depth int) error {
		if depth > maxRecursionDepth {
			return fmt.Errorf("max recursion depth exceeded: %d", maxRecursionDepth)
		}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldPath := append(parentPath, i)

			// check if the column is anonymous struct
			if field.Anonymous && field.Type.Kind() == reflect.Struct {
				if err := collectMeta(field.Type, fieldPath, depth+1); err != nil {
					return err
				}
				continue
			}

			if tag, ok := field.Tag.Lookup("db"); ok {
				if !isValidColumnName(tag) {
					return fmt.Errorf("invalid column: %s", tag)
				}
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
		return nil
	}

	return metas, collectMeta(t, nil, 0)
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

func isValidColumnName(name string) bool {
	if name == "" {
		return false
	}

	// Check for maximum length (for PostgreSQL it's usually 63 characters)
	if len(name) > 63 {
		return false
	}

	// First character must be a letter or underscore
	firstChar := name[0]
	if !((firstChar >= 'a' && firstChar <= 'z') ||
		(firstChar >= 'A' && firstChar <= 'Z') ||
		firstChar == '_') {
		return false
	}

	// Check other characters
	for i := 1; i < len(name); i++ {
		c := name[i]
		// PostgreSQL identifiers allow letters, digits, $ and _
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '_' || c == '$') {
			return false
		}
	}

	return true
}
