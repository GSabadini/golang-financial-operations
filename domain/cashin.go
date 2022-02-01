package domain

import (
	"context"
	"time"
)

type CashInCreator interface {
	Create(context.Context, CashIn) error
}

type CashIn struct {
	id             string
	idempotenceKey string
	accountID      int
	amount         int
	currency       string
	createdAt      time.Time
}

func NewCashIn(
	id string,
	idempotenceKey string,
	accountID int,
	amount int,
	currency string,
	createdAt time.Time,
) CashIn {
	return CashIn{
		id:             id,
		idempotenceKey: idempotenceKey,
		accountID:      accountID,
		amount:         amount,
		currency:       currency,
		createdAt:      createdAt,
	}
}

func (c CashIn) Build(cashin CashIn) {
	return
}
