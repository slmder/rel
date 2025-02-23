package db

type PKStrategy int

const (
	PkStrategySequence PKStrategy = iota + 1
	PkStrategyGenerated
)
