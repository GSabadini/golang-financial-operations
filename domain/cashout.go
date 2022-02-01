package domain

import (
	"context"
	"time"
)

type CashOutCreator interface {
	Create(context.Context, CashOut) error
}

type CashOut struct {
	id             string
	idempotenceKey string
	accountID      int
	amount         int
	currency       string
	createdAt      time.Time
}

func NewCashOut(
	id string,
	idempotenceKey string,
	accountID int,
	amount int,
	currency string,
	createdAt time.Time,
) CashOut {
	return CashOut{
		id:             id,
		idempotenceKey: idempotenceKey,
		accountID:      accountID,
		amount:         amount,
		currency:       currency,
		createdAt:      createdAt,
	}
}
