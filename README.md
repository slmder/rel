[![codecov](https://codecov.io/gh/slmder/rel/graph/badge.svg?token=K7Y8MN8V0H)](https://codecov.io/gh/slmder/rel)
# Rel - Golang data access layer for postgres

`Rel` is a package that provides a basic abstraction layer for working with relational databases in Go. It implements functionality similar to an ORM but uses low-level SQL queries, improving performance and flexibility.

## Features

- Simple and flexible database interaction through structured types.
- SQL queries are generated using reflection and metadata on application startup.
- Support for CRUD operations: `Insert`, `Update`, `Delete`, `Find`, `FindBy`, `FindOneBy`.
- Ability to work with various data types provided via generics.
- Automatic query generation based on data structures.
- Simple sql query builder [qbuilder](qbuilder)

## Installation

To use the package, simply add it to your project with `go get`:

```bash
go get github.com/slmder/rel
```

## Usage

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/slmder/rel"
)

type YourEntity struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func main() {
	// Database connection
	connStr := "user=username dbname=mydb sslmode=disable"
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Create a Relation for your entity
	repository, err := rel.NewRelation[YourEntity]("your_table_name", dbConn)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new entity
	entity := &YourEntity{Name: "Test Name"}

	// Save the entity to the database
	err = repository.Insert(context.Background(), entity)
	if err != nil {
		log.Fatal("Error saving entity:", err)
	}
	fmt.Println("Entity saved, id=" + strconv.Itoa(entity.ID))

	// Find the entity by ID
	foundEntity, err := repository.Find(context.Background(), 1)
	if err != nil {
		log.Fatal("Error finding entity:", err)
	}
	fmt.Printf("Found entity: %+v\n", foundEntity)

	// Update the entity
	entity.Name = "Updated Name"
	err = repository.Update(context.Background(), entity)
	if err != nil {
		log.Fatal("Error updating entity:", err)
	}
	fmt.Println("Entity updated successfully!")

	// Delete the entity
	err = repository.Delete(context.Background(), 1)
	if err != nil {
		log.Fatal("Error deleting entity:", err)
	}
	fmt.Println("Entity deleted successfully!")
}
```

## Notes
Each entity must be a struct type, where fields can be annotated with db tags to specify the corresponding columns in the table.
Cond is a structure that represents a condition for searching in the database. It allows you to build flexible queries.