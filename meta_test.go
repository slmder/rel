package rel

import (
	"reflect"
	"testing"
)

// ColumnMeta.Identifier()
func TestColumnMeta_Identifier(t *testing.T) {
	tests := []struct {
		name     string
		column   ColumnMeta
		expected string
	}{
		{"Simple name", ColumnMeta{name: "username"}, `"username"`},
		{"Name with space", ColumnMeta{name: "user name"}, `"user name"`},
		{"Reserved keyword", ColumnMeta{name: "order"}, `"order"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.column.Identifier(); got != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, got)
			}
		})
	}
}

// ListColumnMeta.Identifiers() и Names()
func TestListColumnMeta_Methods(t *testing.T) {
	columns := ListColumnMeta{
		&ColumnMeta{name: "id"},
		&ColumnMeta{name: "name"},
		&ColumnMeta{name: "email"},
	}

	expectedIdentifiers := []string{`"id"`, `"name"`, `"email"`}
	expectedNames := []string{"id", "name", "email"}

	if got := columns.Identifiers(); !reflect.DeepEqual(got, expectedIdentifiers) {
		t.Errorf("Identifiers() = %v, expected %v", got, expectedIdentifiers)
	}

	if got := columns.Names(); !reflect.DeepEqual(got, expectedNames) {
		t.Errorf("Names() = %v, expected %v", got, expectedNames)
	}
}

// columnsMeta()
func TestColumnsMeta(t *testing.T) {
	type User struct {
		ID    int    `db:"id"`
		Name  string `db:"name"`
		Email string `db:"email"`
	}

	columns, err := columnsMeta[User]("id")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := ListColumnMeta{
		&ColumnMeta{name: "id", pk: true},
		&ColumnMeta{name: "name", pk: false},
		&ColumnMeta{name: "email", pk: false},
	}

	if !equalColumnMetaLists(columns, expected) {
		t.Errorf("Expected %v, got %v", expected, columns)
	}
}

// pkColumns()
func TestPKColumns(t *testing.T) {
	columns := ListColumnMeta{
		&ColumnMeta{name: "id", pk: true},
		&ColumnMeta{name: "name", pk: false},
		&ColumnMeta{name: "email", pk: false},
	}

	expected := ListColumnMeta{
		&ColumnMeta{name: "id", pk: true},
	}

	if got := pkColumns(columns); !reflect.DeepEqual(got, expected) {
		t.Errorf("pkColumns() = %v, expected %v", got, expected)
	}
}

// insertColumns()
func TestInsertColumns(t *testing.T) {
	columns := ListColumnMeta{
		&ColumnMeta{name: "id", pk: true},
		&ColumnMeta{name: "name", pk: false},
		&ColumnMeta{name: "email", pk: false},
	}

	tests := []struct {
		name       string
		pkStrategy PKStrategy
		expected   ListColumnMeta
	}{
		{"With generated PK", PkStrategyGenerated, columns},
		{"Without generated PK", PkStrategySequence, columns[1:]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := insertColumns(columns, tt.pkStrategy); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("insertColumns() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

// updateColumns()
func TestUpdateColumns(t *testing.T) {
	columns := ListColumnMeta{
		&ColumnMeta{name: "id", pk: true},
		&ColumnMeta{name: "name", pk: false},
		&ColumnMeta{name: "email", pk: false},
	}

	expected := ListColumnMeta{
		&ColumnMeta{name: "name", pk: false},
		&ColumnMeta{name: "email", pk: false},
	}

	if got := updateColumns(columns); !reflect.DeepEqual(got, expected) {
		t.Errorf("updateColumns() = %v, expected %v", got, expected)
	}
}

// filter()
func TestFilter(t *testing.T) {
	columns := ListColumnMeta{
		&ColumnMeta{name: "id", pk: true},
		&ColumnMeta{name: "name", pk: false},
		&ColumnMeta{name: "email", pk: false},
	}

	tests := []struct {
		name     string
		filterFn func(*ColumnMeta) bool
		expected ListColumnMeta
	}{
		{
			"Only PK",
			func(cm *ColumnMeta) bool { return cm.pk },
			ListColumnMeta{columns[0]},
		},
		{
			"Non-PK",
			func(cm *ColumnMeta) bool { return !cm.pk },
			ListColumnMeta{columns[1], columns[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filter(columns, tt.filterFn); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("filter() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

// NewMeta()
func TestNewMeta(t *testing.T) {
	type User struct {
		ID    int    `db:"id"`
		Name  string `db:"name"`
		Email string `db:"email"`
	}

	meta, err := NewMeta[User](PkStrategyGenerated, "id")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if meta.PkStrategy() != PkStrategyGenerated {
		t.Errorf("Expected PkStrategy %v, got %v", PkStrategyGenerated, meta.PkStrategy())
	}

	expectedPK := ListColumnMeta{&ColumnMeta{name: "id", pk: true}}

	if !equalColumnMetaLists(meta.PKColumns(), expectedPK) {
		t.Errorf("Expected PK columns %v, got %v", expectedPK, meta.PKColumns())
	}

	expectedInsert := ListColumnMeta{
		&ColumnMeta{name: "id", pk: true},
		&ColumnMeta{name: "name", pk: false},
		&ColumnMeta{name: "email", pk: false},
	}
	if !equalColumnMetaLists(meta.InsertColumns(), expectedInsert) {
		t.Errorf("Expected insert columns %v, got %v", expectedInsert, meta.InsertColumns())
	}
}

func equalColumnMetaLists(a, b ListColumnMeta) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].name != b[i].name || a[i].pk != b[i].pk {
			return false
		}
	}
	return true
}
