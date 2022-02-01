package domain

import (
	"context"
	"time"
)

type TransferCreator interface {
	Create(context.Context, Transfer) error
}

type Transfer struct {
	id             string
	idempotenceKey string
	toAccountID    int
	fromAccountID  int
	amount         int
	currency       string
	createdAt      time.Time
}

func NewTransfer(
	id string,
	idempotenceKey string,
	toAccountID int,
	fromAccountID int,
	amount int,
	currency string,
	createdAt time.Time,
) Transfer {
	return Transfer{
		id: id, idempotenceKey: idempotenceKey,
		toAccountID:   toAccountID,
		fromAccountID: fromAccountID,
		amount:        amount,
		currency:      currency,
		createdAt:     createdAt,
	}
}

func (t Transfer) Build() {

}
