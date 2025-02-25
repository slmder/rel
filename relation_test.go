package db

import (
	"context"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

type TimeStamps struct {
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}

type entitySerialID struct {
	TimeStamps
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

//goland:noinspection SqlNoDataSourceInspection
func TestRelation_InsertSerialID(t *testing.T) {
	tests := []struct {
		name      string
		entity    *entitySerialID
		expectErr bool
	}{
		{
			name:      "successful save",
			entity:    &entitySerialID{ID: 1, Name: "Test Name", TimeStamps: TimeStamps{Created: time.Now(), Updated: time.Now()}},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock db: %v", err)
			}
			eq := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO entities ("created", "updated", "name") VALUES ($1, $2, $3) RETURNING "created", "updated", "id", "name"`))
			eq.WithArgs(tt.entity.Created, tt.entity.Updated, tt.entity.Name)
			eq.WillReturnRows(
				sqlmock.NewRows([]string{"created", "updated", "id", "name"}).
					AddRow(tt.entity.Created, tt.entity.Updated, tt.entity.ID, tt.entity.Name),
			)

			rel, err := NewRelation[entitySerialID]("entities", mockDB)
			if err != nil {
				t.Fatalf("failed to create relation: %v", err)
			}

			var ent = entitySerialID{
				TimeStamps: tt.entity.TimeStamps,
				Name:       tt.entity.Name,
			}
			if err = rel.Insert(context.Background(), &ent); !tt.expectErr && err != nil {
				t.Fatalf("failed to save entitySerialID: %v", err)
			}

			if ent.ID != tt.entity.ID {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}
		})
	}
}

//goland:noinspection SqlNoDataSourceInspection
func TestRelation_UpdateSerialID(t *testing.T) {
	tests := []struct {
		name      string
		entity    *entitySerialID
		expectErr bool
	}{
		{
			name:      "successful update",
			entity:    &entitySerialID{ID: 1, Name: "Test Name", TimeStamps: TimeStamps{Created: time.Now(), Updated: time.Now()}},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock db: %v", err)
			}
			eq := mock.ExpectQuery(regexp.QuoteMeta(`UPDATE entities SET "created" = $1, "updated" = $2, "name" = $3 WHERE "id" = $4 RETURNING "created", "updated", "id", "name"`))
			eq.WithArgs(tt.entity.Created, tt.entity.Updated, tt.entity.Name, tt.entity.ID)
			eq.WillReturnRows(
				sqlmock.NewRows([]string{"created", "updated", "id", "name"}).
					AddRow(tt.entity.Created, tt.entity.Updated, tt.entity.ID, tt.entity.Name),
			)

			rel, err := NewRelation[entitySerialID]("entities", mockDB)
			if err != nil {
				t.Fatalf("failed to create relation: %v", err)
			}

			var ent = entitySerialID{
				ID:         tt.entity.ID,
				TimeStamps: tt.entity.TimeStamps,
				Name:       tt.entity.Name,
			}
			if err = rel.Update(context.Background(), &ent); !tt.expectErr && err != nil {
				t.Fatalf("failed to save entitySerialID: %v", err)
			}

			if ent.ID != tt.entity.ID {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}
		})
	}
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func TestRelation_DeleteSerialID(t *testing.T) {
	tests := []struct {
		name      string
		entity    *entitySerialID
		expectErr bool
	}{
		{
			name:      "successful delete",
			entity:    &entitySerialID{ID: 1, Name: "Test Name", TimeStamps: TimeStamps{Created: time.Now(), Updated: time.Now()}},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock db: %v", err)
			}
			eq := mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM entities WHERE "id" = $1`))
			eq.WithArgs(tt.entity.ID)
			eq.WillReturnResult(sqlmock.NewResult(0, 1))
			rel, err := NewRelation[entitySerialID]("entities", mockDB)
			if err != nil {
				t.Fatalf("failed to create relation: %v", err)
			}

			var ent = entitySerialID{
				ID:         tt.entity.ID,
				TimeStamps: tt.entity.TimeStamps,
				Name:       tt.entity.Name,
			}
			if err = rel.Delete(context.Background(), ent.ID); !tt.expectErr && err != nil {
				t.Fatalf("failed to save entitySerialID: %v", err)
			}

			if ent.ID != tt.entity.ID {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}
		})
	}
}

//func TestRelation_Find(t *testing.T) {
//	tests := []struct {
//		name      string
//		id        interface{}
//		mockDB    func() *sql.DB
//		expected  interface{}
//		expectErr bool
//	}{
//		{
//			name: "successful find",
//			id:   1,
//			mockDB: func() *sql.DB {
//				// Мокаем успешный ответ от базы
//				dbMock := &sql.DB{}
//				return dbMock
//			},
//			expected:  &YourEntity{ID: 1, Name: "Found Name"},
//			expectErr: false,
//		},
//		{
//			name: "find error",
//			id:   1,
//			mockDB: func() *sql.DB {
//				// Мокаем ошибку при поиске
//				dbMock := &sql.DB{}
//				return dbMock
//			},
//			expected:  nil,
//			expectErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			db := tt.mockDB()
//
//			rel, err := db.NewRelation("your_table_name", db)
//			if err != nil {
//				t.Fatalf("failed to create relation: %v", err)
//			}
//
//			result, err := rel.Find(context.Background(), tt.id)
//			if tt.expectErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, tt.expected, result)
//			}
//		})
//	}
//}
//
//func TestRelation_FindBy(t *testing.T) {
//	tests := []struct {
//		name      string
//		cond      db.Cond
//		mockDB    func() *sql.DB
//		expected  []interface{}
//		expectErr bool
//	}{
//		{
//			name: "successful find by condition",
//			cond: db.Cond{"Name": "Test Name"},
//			mockDB: func() *sql.DB {
//				// Мокаем успешный ответ от базы
//				dbMock := &sql.DB{}
//				return dbMock
//			},
//			expected: []interface{}{
//				&YourEntity{ID: 1, Name: "Test Name"},
//			},
//			expectErr: false,
//		},
//		{
//			name: "find by error",
//			cond: db.Cond{"Name": "Test Name"},
//			mockDB: func() *sql.DB {
//				// Мокаем ошибку при поиске
//				dbMock := &sql.DB{}
//				return dbMock
//			},
//			expected:  nil,
//			expectErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			db := tt.mockDB()
//
//			rel, err := db.NewRelation("your_table_name", db)
//			if err != nil {
//				t.Fatalf("failed to create relation: %v", err)
//			}
//
//			result, err := rel.FindBy(context.Background(), tt.cond)
//			if tt.expectErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, tt.expected, result)
//			}
//		})
//	}
//}
