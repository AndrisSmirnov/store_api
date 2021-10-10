package types

type TransactionType uint8

const (
	TransactionType_WIN TransactionType = iota + 1
	TransactionType_ADD
	TransactionType_DEC
)
