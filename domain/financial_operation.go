package domain

import (
	"context"
	"errors"
)

type (
	FinancialOperator interface {
		Process(
			context.Context,
			Authorizer,
			RulesEvaluator,
			Ledger,
			Transactor,
			FinancialOperationStreamer,
		) (FinancialOperation, error)
	}

	Authorizer interface {
		Create(context.Context, AuthorizationInput) (int, error)
		FindByIdempotenceKey(context.Context, AuthorizationInput) (Authorization, error)
	}

	RulesEvaluator interface {
		SpendingControls(context.Context, RulesInput) error
	}

	Ledger interface {
		AvailableEntries(context.Context, LedgerInput) error
	}

	Acquirer interface {
		Sale(context.Context, string, string, string, string) error
	}

	Transactor interface {
		Transaction(context.Context, TransactionInput) (int, error)
	}

	FinancialOperationStreamer interface {
		Publish(context.Context, StreamInput)
	}
)

var (
	ErrDuplicatedFinancialOperation = errors.New("duplicated financial operation")
)

type FinancialOperation struct {
	id             string
	idempotenceKey string
	cashIn         CashIn
	cashOut        CashOut
	transfer       Transfer
}

func (f FinancialOperation) ID() string {
	return f.id
}

func (f FinancialOperation) IdempotenceKey() string {
	return f.idempotenceKey
}

func (f FinancialOperation) AlreadyExists() bool {
	return f.id != ""
}
