package rel

import (
	"errors"
	"reflect"
	"testing"
)

type mockScanner struct {
	values []any
	err    error
}

func (m *mockScanner) scan(dest ...any) error {
	if m.err != nil {
		return m.err
	}
	for i := range dest {
		reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(m.values[i]))
	}
	return nil
}

type TestEntity struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func TestScanRow(t *testing.T) {
	tests := []struct {
		name    string
		values  []any
		scanErr error
		wantErr bool
	}{
		{
			name:    "Successful scan",
			values:  []any{1, "John"},
			wantErr: false,
		},
		{
			name:    "Scan with error",
			scanErr: errors.New("scan error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockScanner{values: tt.values, err: tt.scanErr}
			meta, err := NewMeta[TestEntity](PkStrategyGenerated, "id")
			if err != nil {
				t.Fatalf("meta error: %v", err)
			}
			var entity TestEntity
			err = scanRow[TestEntity](mock.scan, meta, &entity)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetFieldsPointers(t *testing.T) {
	tests := []struct {
		name    string
		fields  []string
		wantErr bool
	}{
		{
			name:    "Valid fields",
			fields:  []string{"id", "name"},
			wantErr: false,
		},
		{
			name:    "Invalid field",
			fields:  []string{"Unknown"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, err := NewMeta[TestEntity](PkStrategyGenerated, "id")
			if err != nil {
				t.Fatalf("meta error: %v", err)
			}
			entity := TestEntity{}
			_, err = getFieldsPointers[TestEntity](tt.fields, meta, &entity)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
