package qbuilder

import "strings"

type SortDirection int

const (
	SortDirectionASC SortDirection = iota
	SortDirectionDESC
)

func (d SortDirection) String() string {
	return [...]string{"ASC", "DESC"}[d]
}

func SortDirFromString(d string) SortDirection {
	if strings.ToLower(d) == SortDirectionDESC.String() {
		return SortDirectionDESC
	}
	return SortDirectionASC
}

type rowLevelLockMode int

const (
	LockModeUpdate rowLevelLockMode = iota
	LockModeUpdateNowait
	LockModeShare
	LockModeShareNowait
	LockModeNoKeyUpdate
	LockModeKeyShare
	LockModeUpdateSkipLocked
)

func (m rowLevelLockMode) String() string {
	return [...]string{
		"UPDATE", "UPDATE NOWAIT", "SHARE", "SHARE NOWAIT", "NO KEY UPDATE", "KEY SHARE", "UPDATE SKIP LOCKED",
	}[m]
}

type joinType string

const (
	JoinTypeLeft  = joinType("LEFT")
	joinTypeRight = joinType("RIGHT")
	joinTypeInner = joinType("INNER")
	joinTypeCross = joinType("CROSS")
)

func (d joinType) String() string {
	return string(d)
}
