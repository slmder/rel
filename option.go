package db

type Option[T any] func(*Relation[T])

func PK[T any](pk ...string) Option[T] {
	return func(db *Relation[T]) {
		db.pk = pk
	}
}

func PKStrategyGenerated[T any](db *Relation[T]) {
	db.pkStrategy = PkStrategyGenerated
}
