package domain

import (
	"errors"
)

type FinancialOperationBuilder interface {
	Build(string) (Payment, error)
}

var (
	ErrDuplicatedPayment = errors.New("duplicated payment")
)

type Payment struct {
	id             string
	idempotenceKey string
}

func (p Payment) ID() string {
	return p.id
}

func (p Payment) IdempotenceKey() string {
	return p.idempotenceKey
}

func (p Payment) AlreadyExists() bool {
	return p.id != ""
}
