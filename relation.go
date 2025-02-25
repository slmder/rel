package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"

	"github.com/slmder/rel/qbuilder"
)

const defaultPkColumn = "id"

// Relation is a database abstraction layer that provides basic CRUD operations for a given type.
// Uses reflection prebuild queries for the given type and table
type Relation[T any] struct {
	*sql.DB
	// metadata
	M *Metadata[T]

	// Relation name
	name string
	// primary key columns
	pk []string
	// primary key strategy
	pkStrategy PKStrategy
	// prebuilt queries
	// getOneQ is a prebuilt query to get a single entitySerialID by primary key
	getOneQ string
	// insertQ is a prebuilt query to insert an entitySerialID
	insertQ string
	// updateQ is a prebuilt query to update an entitySerialID
	updateQ string
	// deleteQ is a prebuilt query to delete an entitySerialID
	deleteQ string
	// findByQ is a prebuilt query to find entities by operator
	findByQ qbuilder.SelectBuilder
}

// NewRelation creates a new Relation instance for the given type and table.
// Relation requires a primary key to be specified at least one column (by default it is 'id').
func NewRelation[T any](name string, db *sql.DB, opts ...Option[T]) (*Relation[T], error) {
	d := &Relation[T]{
		name: name,
		DB:   db,
		pk:   []string{defaultPkColumn},
	}
	for _, o := range opts {
		o(d)
	}
	var err error
	d.M, err = NewMeta[T](d.pkStrategy, d.pk...)
	if err != nil {
		return nil, fmt.Errorf("create '%s' meta: %w", name, err)
	}

	d.insertQ = buildInsertQuery(name, d.M)
	d.updateQ = buildUpdateQuery(name, d.M)
	d.deleteQ = buildDeleteQuery(name, d.M)
	d.getOneQ = buildGetOneQuery(name, d.M)
	d.findByQ = buildFindByQuery(name, d.M)

	return d, nil
}

// Rel returns the table name
func (r *Relation[T]) Rel() string {
	return r.name
}

// InsertArgsFrom returns arguments for insert for the given entitySerialID
func (r *Relation[T]) InsertArgsFrom(e *T) []any {
	return getFieldsValues(r.M.InsertColumns().Names(), r.M, e)
}

// UpdateArgsFrom returns arguments for update for the given entitySerialID
func (r *Relation[T]) UpdateArgsFrom(e *T) []any {
	return getFieldsValues(r.M.UpdateColumns().Names(), r.M, e)
}

// Insert inserts an entitySerialID
func (r *Relation[T]) Insert(ctx context.Context, entity *T) error {
	args := getFieldsValues(r.M.InsertColumns().Names(), r.M, entity)
	row := r.DB.QueryRowContext(ctx, r.insertQ, args...)

	return scanRow(row.Scan, r.M, entity)
}

// Update updates an entitySerialID
func (r *Relation[T]) Update(ctx context.Context, entity *T) error {
	args := getFieldsValues(r.M.UpdateColumns().Names(), r.M, entity)
	args = append(args, getFieldsValues(r.M.PKColumns().Names(), r.M, entity)...)
	row := r.DB.QueryRowContext(ctx, r.updateQ, args...)

	return scanRow(row.Scan, r.M, entity)
}

// Delete deletes an entitySerialID by given id
func (r *Relation[T]) Delete(ctx context.Context, id ...any) error {
	_, err := r.DB.ExecContext(ctx, r.deleteQ, id...)

	if err != nil {
		return fmt.Errorf("delete record: %w", err)
	}
	return nil
}

// Find finds single entitySerialID by given id
func (r *Relation[T]) Find(ctx context.Context, id ...any) (T, error) {
	var entity T
	row := r.DB.QueryRowContext(ctx, r.getOneQ, id...)

	return entity, r.Scan(row.Scan, &entity)
}

// FindBy finds all entities by given operator
func (r *Relation[T]) FindBy(ctx context.Context, cond Cond) ([]T, error) {
	var items []T
	query := r.findByQ.Copy()
	args, expr := cond.Split()

	if len(expr) > 0 {
		for _, e := range expr {
			query.AndWhere(e)
		}
	}

	rows, err := r.DB.QueryContext(ctx, query.ToSQL(), args...)
	if err != nil {
		return nil, fmt.Errorf("db find by query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entity T
		if err := scanRow(rows.Scan, r.M, &entity); err != nil {
			return nil, fmt.Errorf("db find by scan: %w", err)
		}
		items = append(items, entity)
	}

	return items, nil
}

// FindOneBy finds single entitySerialID by given operator
func (r *Relation[T]) FindOneBy(ctx context.Context, cond Cond) (T, error) {
	var entity T
	query := r.findByQ.Copy()
	args, expr := cond.Split()

	if len(expr) > 0 {
		for _, e := range expr {
			query.AndWhere(e)
		}
	}
	query.Limit(1)
	row := r.DB.QueryRowContext(ctx, query.ToSQL(), args...)

	return entity, r.Scan(row.Scan, &entity)
}

// Scan scans a single row into the entitySerialID
func (r *Relation[T]) Scan(sf scanFunc, dst *T) error {
	return scanRow(sf, r.M, dst)
}

// ArgsAdd adds an argument to the slice and returns a placeholder for it.
func ArgsAdd(args []any, arg any) string {
	args = append(args, arg)
	return fmt.Sprintf("$%d", len(args))
}

// buildInsertQuery prebuilds a query to insert an entitySerialID
func buildInsertQuery[T any](rel string, m *Metadata[T]) string {
	qb := qbuilder.Insert(rel)
	qb.Columns(m.InsertColumns().Identifiers()...)
	qb.Values(getArgsPlaceholders(len(m.InsertColumns())))
	qb.Returning(m.Columns().Identifiers()...)

	return qb.ToSQL()
}

// buildUpdateQuery prebuilds a query to update an entitySerialID
func buildUpdateQuery[T any](rel string, m *Metadata[T]) string {
	qb := qbuilder.Update(rel)

	var i int
	for _, col := range m.UpdateColumns() {
		qb.Set(col.Identifier(), "$"+strconv.Itoa(i+1))
		i++
	}

	for _, col := range m.PKColumns() {
		qb.AndWhere(col.Identifier() + " = $" + strconv.Itoa(i+1))
		i++
	}
	qb.Returning(m.Columns().Identifiers()...)

	return qb.ToSQL()
}

// buildDeleteQuery prebuilds a query to delete an entitySerialID
func buildDeleteQuery[T any](rel string, m *Metadata[T]) string {
	qb := qbuilder.Delete(rel)

	for i, col := range m.PKColumns() {
		qb.AndWhere(col.Identifier() + " = $" + strconv.Itoa(i+1))
	}

	return qb.ToSQL()
}

// buildGetOneQuery prebuilds a query to get a single entitySerialID
func buildGetOneQuery[T any](rel string, m *Metadata[T]) string {
	qb := qbuilder.Select(m.Columns().Identifiers()...)
	qb.From(rel)

	for _, col := range m.PKColumns() {
		qb.Where(col.Identifier(), col.name)
	}

	return qb.Limit(1).ToSQL()
}

// buildFindByQuery prebuilds a query to find many entities
func buildFindByQuery[T any](rel string, m *Metadata[T]) qbuilder.SelectBuilder {
	qb := qbuilder.Select(m.Columns().Identifiers()...)
	qb.From(rel)

	return qb.Copy()
}

// getFieldsValues returns values of the fields of the struct.
func getFieldsValues[T any](fields []string, meta *Metadata[T], entity *T) []any {
	v := reflect.ValueOf(entity)

	// If a pointer is passed, we get the value it points to.
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("entitySerialID must be a struct or pointer to a struct")
	}

	result := make([]any, len(fields))
	// Retrieving column values.
	for i, name := range fields {
		if fieldInfo, found := meta.columnsMap[name]; found {
			fieldValue := v
			for _, index := range fieldInfo.path {
				fieldValue = fieldValue.Field(index)
			}
			result[i] = fieldValue.Interface()
		} else {
			panic(fmt.Sprintf("column not found: %s", name))
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
