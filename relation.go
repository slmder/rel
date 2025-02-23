package db

import (
	"context"
	"database/sql"
	"fmt"

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
	// primary key Columns
	pk []string
	// primary key strategy
	pkStrategy PKStrategy
	// prebuilt queries
	// getOneQ is a prebuilt query to get a single entity by primary key
	getOneQ string
	// insertQ is a prebuilt query to insert an entity
	insertQ string
	// updateQ is a prebuilt query to update an entity
	updateQ string
	// deleteQ is a prebuilt query to delete an entity
	deleteQ string
	// findByQ is a prebuilt query to find entities by operator
	findByQ string
}

// NewRelation creates a new Relation instance for the given type and table.
// Relation requires a primary key to be specified at least one column.
// Spread operator is used to allow composite keys and simplify the usage.
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

// InsertArgsFrom returns arguments for the given entity
func (r *Relation[T]) InsertArgsFrom(e *T) []any {
	return getFieldsValues(r.M.InsertColumns.Names(), r.M, e)
}

// Insert returns a prebuilt query to insert an entity.
// It's useful when you need upsert an entity.
func (r *Relation[T]) Insert() *qbuilder.InsertBuilder {
	return qbuilder.Insert(r.name).
		Columns(r.M.InsertColumns.Identifiers()...).
		Values(getArgsPlaceholders(len(r.M.InsertColumns))).
		Returning(r.M.Columns.Identifiers()...)
}

// Save inserts an entity
func (r *Relation[T]) Save(ctx context.Context, entity *T) error {
	args := getFieldsValues(r.M.InsertColumns.Names(), r.M, entity)
	row := r.DB.QueryRowContext(ctx, r.updateQ, args...)

	return scanPK(row.Scan, r.M, entity)
}

// Change updates an entity
func (r *Relation[T]) Change(ctx context.Context, entity *T) error {
	args := getFieldsValues(r.M.UpdateColumns.Names(), r.M, entity)
	_, err := r.DB.ExecContext(ctx, r.updateQ, args...)

	if err != nil {
		return fmt.Errorf("update record: %w", err)
	}
	return nil
}

// Remove deletes an entity by given id
func (r *Relation[T]) Remove(ctx context.Context, id ...any) error {
	_, err := r.DB.ExecContext(ctx, r.deleteQ, id...)

	if err != nil {
		return fmt.Errorf("delete record: %w", err)
	}
	return nil
}

// Find finds single entity by given id
func (r *Relation[T]) Find(ctx context.Context, id ...any) (T, error) {
	var entity T
	row := r.DB.QueryRowContext(ctx, r.getOneQ, id...)

	return entity, r.Scan(row.Scan, &entity)
}

// FindBy finds all entities by given operator
func (r *Relation[T]) FindBy(ctx context.Context, cond Cond) ([]T, error) {
	var items []T
	query := r.findByQ
	args, expr := cond.Split()

	if len(expr) > 0 {
		query += " WHERE " + expr[0]
	}

	if len(expr) > 1 {
		for _, e := range expr[1:] {
			query += " AND " + e
		}
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
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

// FindOneBy finds single entity by given operator
func (r *Relation[T]) FindOneBy(ctx context.Context, cond Cond) (T, error) {
	var entity T
	query := r.findByQ
	args, expr := cond.Split()
	if len(expr) > 0 {
		query += " WHERE " + expr[0]
	}
	if len(expr) > 1 {
		for _, e := range expr[1:] {
			query += " AND " + e
		}
	}
	query += " LIMIT 1;"
	row := r.DB.QueryRowContext(ctx, query, args...)

	return entity, r.Scan(row.Scan, &entity)
}

// Scan scans a single row into the entity
func (r *Relation[T]) Scan(sf scanFunc, dst *T) error {
	return scanRow(sf, r.M, dst)
}

// buildInsertQuery prebuilds a query to insert an entity
func buildInsertQuery[T any](rel string, m *Metadata[T]) string {
	qb := qbuilder.Insert(rel)
	qb.Columns(m.InsertColumns.Identifiers()...)
	qb.Values(getArgsPlaceholders(len(m.InsertColumns)))
	qb.Returning(m.Columns.Identifiers()...)

	return qb.ToSQL()
}

// buildUpdateQuery prebuilds a query to update an entity
func buildUpdateQuery[T any](rel string, m *Metadata[T]) string {
	qb := qbuilder.Update(rel)

	for _, col := range m.UpdateColumns {
		qb.Set(col.Identifier(), col.name)
	}

	for _, col := range m.PKColumns {
		qb.AndWhere(col.Identifier(), col.name)
	}

	return qb.ToSQL()
}

// buildDeleteQuery prebuilds a query to delete an entity
func buildDeleteQuery[T any](rel string, m *Metadata[T]) string {
	qb := qbuilder.Delete(rel)

	for _, col := range m.PKColumns {
		qb.AndWhere(col.Identifier(), col.name)
	}

	return qb.ToSQL()
}

// buildGetOneQuery prebuilds a query to get a single entity
func buildGetOneQuery[T any](rel string, m *Metadata[T]) string {
	qb := qbuilder.Select(m.Columns.Identifiers()...)
	qb.From(rel)

	for _, col := range m.PKColumns {
		qb.Where(col.Identifier(), col.name)
	}

	return qb.Limit(1).ToSQL()
}

// buildFindByQuery prebuilds a query to find many entities
func buildFindByQuery[T any](rel string, m *Metadata[T]) string {
	qb := qbuilder.Select(m.Columns.Identifiers()...)
	qb.From(rel)

	return qb.ToSQL()
}
