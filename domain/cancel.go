package domain

import (
	"context"
	"time"
)

type CancelCreator interface {
	Create(context.Context, Cancel) error
}

type Cancel struct {
	id             string
	idempotenceKey string
	accountID      string
	paymentID      string
	amount         int
	currency       string
	createdAt      time.Time
}

func NewCancel(
	id string,
	idempotenceKey string,
	accountID string,
	paymentID string,
	amount int,
	currency string,
	createdAt time.Time,
) Cancel {
	return Cancel{
		id:             id,
		idempotenceKey: idempotenceKey,
		accountID:      accountID,
		paymentID:      paymentID,
		amount:         amount,
		currency:       currency,
		createdAt:      createdAt,
	}
}
