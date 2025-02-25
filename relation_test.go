package rel

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
func TestRelationSerial_InsertSerialID(t *testing.T) {
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
			eq := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "entities" ("created", "updated", "name") VALUES ($1, $2, $3) RETURNING "created", "updated", "id", "name"`))
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
func TestRelationSerial_UpdateSerialID(t *testing.T) {
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
			eq := mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "entities" SET "created" = $1, "updated" = $2, "name" = $3 WHERE "id" = $4 RETURNING "created", "updated", "id", "name"`))
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
func TestRelationSerial_DeleteSerialID(t *testing.T) {
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
			eq := mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "entities" WHERE "id" = $1`))
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

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func TestRelationSerial_Find(t *testing.T) {
	tests := []struct {
		name      string
		entity    *entitySerialID
		expectErr bool
	}{
		{
			name:      "successful find",
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
			eq := mock.ExpectQuery(regexp.QuoteMeta(`SELECT "created", "updated", "id", "name" FROM "entities" WHERE "id" = $1`))
			eq.WithArgs(tt.entity.ID)
			eq.WillReturnRows(
				sqlmock.NewRows([]string{"created", "updated", "id", "name"}).
					AddRow(tt.entity.Created, tt.entity.Updated, tt.entity.ID, tt.entity.Name),
			)
			rel, err := NewRelation[entitySerialID]("entities", mockDB)
			if err != nil {
				t.Fatalf("failed to create relation: %v", err)
			}

			ent, err := rel.Find(context.Background(), tt.entity.ID)
			if !tt.expectErr && err != nil {
				t.Fatalf("failed to save entitySerialID: %v", err)
			}

			if ent.ID != tt.entity.ID {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}

			if ent.Created != tt.entity.Created {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}

			if ent.Updated != tt.entity.Updated {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}

			if ent.Name != tt.entity.Name {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}

		})
	}
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func TestRelationSerial_FindOneBy(t *testing.T) {
	tests := []struct {
		name      string
		entity    *entitySerialID
		expectErr bool
	}{
		{
			name:      "successful find one by",
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
			eq := mock.ExpectQuery(regexp.QuoteMeta(`SELECT "created", "updated", "id", "name" FROM "entities" WHERE "id" = $1`))
			eq.WithArgs(tt.entity.ID)
			eq.WillReturnRows(
				sqlmock.NewRows([]string{"created", "updated", "id", "name"}).
					AddRow(tt.entity.Created, tt.entity.Updated, tt.entity.ID, tt.entity.Name),
			)
			rel, err := NewRelation[entitySerialID]("entities", mockDB)
			if err != nil {
				t.Fatalf("failed to create relation: %v", err)
			}

			ent, err := rel.FindOneBy(context.Background(), Cond{
				Eq("id", tt.entity.ID),
			})
			if !tt.expectErr && err != nil {
				t.Fatalf("failed to save entitySerialID: %v", err)
			}

			if ent.ID != tt.entity.ID {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}

			if ent.Created != tt.entity.Created {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}

			if ent.Updated != tt.entity.Updated {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}

			if ent.Name != tt.entity.Name {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}

		})
	}
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func TestRelationSerial_FindBy(t *testing.T) {
	tests := []struct {
		name      string
		entity    *entitySerialID
		expectErr bool
	}{
		{
			name:      "successful find by",
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
			eq := mock.ExpectQuery(regexp.QuoteMeta(`SELECT "created", "updated", "id", "name" FROM "entities" WHERE "name" = $1 AND "created" = $2`))
			eq.WithArgs(tt.entity.Name, tt.entity.Created)
			eq.WillReturnRows(
				sqlmock.NewRows([]string{"created", "updated", "id", "name"}).
					AddRow(tt.entity.Created, tt.entity.Updated, tt.entity.ID, tt.entity.Name),
			)
			rel, err := NewRelation[entitySerialID]("entities", mockDB)
			if err != nil {
				t.Fatalf("failed to create relation: %v", err)
			}

			ents, err := rel.FindBy(context.Background(), Cond{
				Eq("name", tt.entity.Name),
				Eq("created", tt.entity.Created),
			})
			if !tt.expectErr && err != nil {
				t.Fatalf("failed to save entitySerialID: %v", err)
			}
			if len(ents) == 0 {
				t.Fatalf("unexpected result")
			}
			ent := ents[0]
			if ent.ID != tt.entity.ID {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}

			if ent.Created != tt.entity.Created {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}

			if ent.Updated != tt.entity.Updated {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}

			if ent.Name != tt.entity.Name {
				t.Fatalf("unexpected ID: %d", ent.ID)
			}

		})
	}
}

type entityCompositeID struct {
	TimeStamps
	IDA  int64  `db:"id_a"`
	IDB  int64  `db:"id_b"`
	Name string `db:"name"`
}

//goland:noinspection SqlNoDataSourceInspection
func TestRelationComposite_InsertSerialID(t *testing.T) {
	tests := []struct {
		name      string
		entity    *entityCompositeID
		expectErr bool
	}{
		{
			name:      "successful save",
			entity:    &entityCompositeID{IDA: 1, IDB: 2, Name: "Test Name", TimeStamps: TimeStamps{Created: time.Now(), Updated: time.Now()}},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock db: %v", err)
			}
			eq := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "entities" ("created", "updated", "id_a", "id_b", "name") VALUES ($1, $2, $3, $4, $5) RETURNING "created", "updated", "id_a", "id_b", "name"`))
			eq.WithArgs(tt.entity.Created, tt.entity.Updated, tt.entity.IDA, tt.entity.IDB, tt.entity.Name)
			eq.WillReturnRows(
				sqlmock.NewRows([]string{"created", "updated", "id_a", "id_b", "name"}).
					AddRow(tt.entity.Created, tt.entity.Updated, tt.entity.IDA, tt.entity.IDB, tt.entity.Name),
			)

			rel, err := NewRelation[entityCompositeID]("entities", mockDB, PKStrategyGenerated, PK[entityCompositeID]("id_a", "id_b"))
			if err != nil {
				t.Fatalf("failed to create relation: %v", err)
			}

			var ent = entityCompositeID{
				IDA:        tt.entity.IDA,
				IDB:        tt.entity.IDB,
				TimeStamps: tt.entity.TimeStamps,
				Name:       tt.entity.Name,
			}
			if err = rel.Insert(context.Background(), &ent); !tt.expectErr && err != nil {
				t.Fatalf("failed to save entityCompositeID: %v", err)
			}

			if ent.IDA != tt.entity.IDA {
				t.Fatalf("unexpected ID: %d", ent.IDA)
			}

			if ent.IDB != tt.entity.IDB {
				t.Fatalf("unexpected ID: %d", ent.IDA)
			}
		})
	}
}

//goland:noinspection SqlNoDataSourceInspection
func TestRelationComposite_UpdateSerialID(t *testing.T) {
	tests := []struct {
		name      string
		entity    *entityCompositeID
		expectErr bool
	}{
		{
			name:      "successful update",
			entity:    &entityCompositeID{IDA: 1, IDB: 2, Name: "Test Name", TimeStamps: TimeStamps{Created: time.Now(), Updated: time.Now()}},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock db: %v", err)
			}
			eq := mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "entities" SET "created" = $1, "updated" = $2, "name" = $3 WHERE "id_a" = $4 AND "id_b" = $5 RETURNING "created", "updated", "id_a", "id_b", "name"`))
			eq.WithArgs(tt.entity.Created, tt.entity.Updated, tt.entity.Name, tt.entity.IDA, tt.entity.IDB)
			eq.WillReturnRows(
				sqlmock.NewRows([]string{"created", "updated", "id_a", "id_b", "name"}).
					AddRow(tt.entity.Created, tt.entity.Updated, tt.entity.IDA, tt.entity.IDB, tt.entity.Name),
			)

			rel, err := NewRelation[entityCompositeID]("entities", mockDB, PKStrategyGenerated, PK[entityCompositeID]("id_a", "id_b"))
			if err != nil {
				t.Fatalf("failed to create relation: %v", err)
			}

			var ent = entityCompositeID{
				IDA:        tt.entity.IDA,
				IDB:        tt.entity.IDB,
				TimeStamps: tt.entity.TimeStamps,
				Name:       tt.entity.Name,
			}
			if err = rel.Update(context.Background(), &ent); !tt.expectErr && err != nil {
				t.Fatalf("failed to save entityCompositeID: %v", err)
			}

			if ent.IDA != tt.entity.IDA {
				t.Fatalf("unexpected ID: %d", ent.IDA)
			}

			if ent.IDB != tt.entity.IDB {
				t.Fatalf("unexpected ID: %d", ent.IDA)
			}
		})
	}
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func TestRelationComposite_DeleteSerialID(t *testing.T) {
	tests := []struct {
		name      string
		entity    *entityCompositeID
		expectErr bool
	}{
		{
			name:      "successful delete",
			entity:    &entityCompositeID{IDA: 1, IDB: 2, Name: "Test Name", TimeStamps: TimeStamps{Created: time.Now(), Updated: time.Now()}},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock db: %v", err)
			}
			eq := mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "entities" WHERE "id_a" = $1 AND "id_b" = $2`))
			eq.WithArgs(tt.entity.IDA, tt.entity.IDB)
			eq.WillReturnResult(sqlmock.NewResult(0, 1))
			rel, err := NewRelation[entityCompositeID]("entities", mockDB, PKStrategyGenerated, PK[entityCompositeID]("id_a", "id_b"))
			if err != nil {
				t.Fatalf("failed to create relation: %v", err)
			}

			var ent = entityCompositeID{
				IDA:        tt.entity.IDA,
				IDB:        tt.entity.IDB,
				TimeStamps: tt.entity.TimeStamps,
				Name:       tt.entity.Name,
			}
			if err = rel.Delete(context.Background(), ent.IDA, ent.IDB); !tt.expectErr && err != nil {
				t.Fatalf("failed to save entitySerialID: %v", err)
			}
		})
	}
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func TestRelationComposite_Find(t *testing.T) {
	tests := []struct {
		name      string
		entity    *entityCompositeID
		expectErr bool
	}{
		{
			name:      "successful find",
			entity:    &entityCompositeID{IDA: 1, IDB: 2, Name: "Test Name", TimeStamps: TimeStamps{Created: time.Now(), Updated: time.Now()}},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock db: %v", err)
			}
			eq := mock.ExpectQuery(regexp.QuoteMeta(`SELECT "created", "updated", "id_a", "id_b", "name" FROM "entities" WHERE "id_a" = $1 AND "id_b" = $2 LIMIT 1`))
			eq.WithArgs(tt.entity.IDA, tt.entity.IDB)
			eq.WillReturnRows(
				sqlmock.NewRows([]string{"created", "updated", "id_a", "id_b", "name"}).
					AddRow(tt.entity.Created, tt.entity.Updated, tt.entity.IDA, tt.entity.IDB, tt.entity.Name),
			)
			rel, err := NewRelation[entityCompositeID]("entities", mockDB, PKStrategyGenerated, PK[entityCompositeID]("id_a", "id_b"))
			if err != nil {
				t.Fatalf("failed to create relation: %v", err)
			}

			ent, err := rel.Find(context.Background(), tt.entity.IDA, tt.entity.IDB)
			if !tt.expectErr && err != nil {
				t.Fatalf("failed to save entityCompositeID: %v", err)
			}

			if ent.IDA != tt.entity.IDA {
				t.Fatalf("unexpected ID: %d", ent.IDA)
			}

			if ent.IDB != tt.entity.IDB {
				t.Fatalf("unexpected ID: %d", ent.IDA)
			}

			if ent.Created != tt.entity.Created {
				t.Fatalf("unexpected Created: %v", ent.Created)
			}

			if ent.Updated != tt.entity.Updated {
				t.Fatalf("unexpected Updated: %v", ent.Updated)
			}

			if ent.Name != tt.entity.Name {
				t.Fatalf("unexpected Name: %s", ent.Name)
			}

		})
	}
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func TestRelationComposite_FindOneBy(t *testing.T) {
	tests := []struct {
		name      string
		entity    *entityCompositeID
		expectErr bool
	}{
		{
			name:      "successful find one by",
			entity:    &entityCompositeID{IDA: 1, IDB: 2, Name: "Test Name", TimeStamps: TimeStamps{Created: time.Now(), Updated: time.Now()}},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock db: %v", err)
			}
			eq := mock.ExpectQuery(regexp.QuoteMeta(`SELECT "created", "updated", "id_a", "id_b", "name" FROM "entities" WHERE "id_a" = $1 AND "id_b" = $2 LIMIT 1`))
			eq.WithArgs(tt.entity.IDA, tt.entity.IDB)
			eq.WillReturnRows(
				sqlmock.NewRows([]string{"created", "updated", "id_a", "id_b", "name"}).
					AddRow(tt.entity.Created, tt.entity.Updated, tt.entity.IDA, tt.entity.IDB, tt.entity.Name),
			)
			rel, err := NewRelation[entityCompositeID]("entities", mockDB, PKStrategyGenerated, PK[entityCompositeID]("id_a", "id_b"))
			if err != nil {
				t.Fatalf("failed to create relation: %v", err)
			}

			ent, err := rel.FindOneBy(context.Background(), Cond{
				Eq("id_a", tt.entity.IDA),
				Eq("id_b", tt.entity.IDB),
			})
			if !tt.expectErr && err != nil {
				t.Fatalf("failed to save entityCompositeID: %v", err)
			}

			if ent.IDA != tt.entity.IDA {
				t.Fatalf("unexpected ID: %d", ent.IDA)
			}

			if ent.IDB != tt.entity.IDB {
				t.Fatalf("unexpected ID: %d", ent.IDA)
			}

			if ent.Created != tt.entity.Created {
				t.Fatalf("unexpected Created: %v", ent.Created)
			}

			if ent.Updated != tt.entity.Updated {
				t.Fatalf("unexpected Updated: %v", ent.Updated)
			}

			if ent.Name != tt.entity.Name {
				t.Fatalf("unexpected Name: %s", ent.Name)
			}

		})
	}
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func TestRelationComposite_FindBy(t *testing.T) {
	tests := []struct {
		name      string
		entity    *entityCompositeID
		expectErr bool
	}{
		{
			name:      "successful find by",
			entity:    &entityCompositeID{IDA: 1, IDB: 2, Name: "Test Name", TimeStamps: TimeStamps{Created: time.Now(), Updated: time.Now()}},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock db: %v", err)
			}
			eq := mock.ExpectQuery(regexp.QuoteMeta(`SELECT "created", "updated", "id_a", "id_b", "name" FROM "entities" WHERE "name" = $1 AND "created" = $2`))
			eq.WithArgs(tt.entity.Name, tt.entity.Created)
			eq.WillReturnRows(
				sqlmock.NewRows([]string{"created", "updated", "id_a", "id_b", "name"}).
					AddRow(tt.entity.Created, tt.entity.Updated, tt.entity.IDA, tt.entity.IDB, tt.entity.Name),
			)
			rel, err := NewRelation[entityCompositeID]("entities", mockDB, PKStrategyGenerated, PK[entityCompositeID]("id_a", "id_b"))
			if err != nil {
				t.Fatalf("failed to create relation: %v", err)
			}

			ents, err := rel.FindBy(context.Background(), Cond{
				Eq("name", tt.entity.Name),
				Eq("created", tt.entity.Created),
			})
			if !tt.expectErr && err != nil {
				t.Fatalf("failed to save entityCompositeID: %v", err)
			}
			if len(ents) == 0 {
				t.Fatalf("unexpected result")
			}
			ent := ents[0]

			if ent.IDA != tt.entity.IDA {
				t.Fatalf("unexpected ID: %d", ent.IDA)
			}

			if ent.IDB != tt.entity.IDB {
				t.Fatalf("unexpected ID: %d", ent.IDA)
			}

			if ent.Created != tt.entity.Created {
				t.Fatalf("unexpected Created: %v", ent.Created)
			}

			if ent.Updated != tt.entity.Updated {
				t.Fatalf("unexpected Updated: %v", ent.Updated)
			}

			if ent.Name != tt.entity.Name {
				t.Fatalf("unexpected Name: %s", ent.Name)
			}
		})
	}
}
