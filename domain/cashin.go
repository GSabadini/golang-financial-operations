package domain

import (
	"context"
)

type CashIn struct {
	idempotenceKey string
	accountID      int
	amount         int
	currency       string
}

func NewCashIn(
	idempotenceKey string,
	accountID int,
	amount int,
	currency string,
) CashIn {
	return CashIn{
		idempotenceKey: idempotenceKey,
		accountID:      accountID,
		amount:         amount,
		currency:       currency,
	}
}

func (c CashIn) Process(
	ctx context.Context,
	authorizer Authorizer,
	rulesEvaluator RulesEvaluator,
	ledger Ledger,
	transactor Transactor,
	streamer FinancialOperationStreamer,
) (FinancialOperation, error) {
	// First step - Authorization
	authorization, err := authorizer.FindByIdempotenceKey(ctx, c.NewAuthorizationInput())
	if err != nil {
		return FinancialOperation{}, err
	}

	if authorization.AlreadyExists() {
		return FinancialOperation{
			id:             authorization.ID(),
			idempotenceKey: authorization.IdempotenceKey(),
		}, ErrDuplicatedAuthorization
	}

	if err := rulesEvaluator.SpendingControls(ctx, c.NewRulesInput()); err != nil {
		return FinancialOperation{}, err
	}

	if err := ledger.AvailableEntries(ctx, c.NewLedgerInput()); err != nil {
		return FinancialOperation{}, err
	}

	authorizationID, err := authorizer.Create(ctx, c.NewAuthorizationInput())
	if err != nil {
		return FinancialOperation{}, err
	}

	// Second step - Transaction
	transactionID, err := transactor.Transaction(ctx, c.NewTransactionInput(authorizationID))
	if err != nil {
		return FinancialOperation{}, err
	}

	// Data stream step
	streamer.Publish(ctx, c.NewStreamInput(authorizationID, transactionID))

	return FinancialOperation{
		id:             authorization.ID(),
		idempotenceKey: authorization.IdempotenceKey(),
	}, nil
}

func (c CashIn) NewLedgerInput() LedgerInput {
	return LedgerInput{
		operation:      "CREDIT",
		currency:       c.currency,
		amount:         c.amount,
		accountID:      c.accountID,
		idempotenceKey: c.idempotenceKey,
	}
}

func (c CashIn) NewRulesInput() RulesInput {
	return RulesInput{
		operation:      "CREDIT",
		currency:       c.currency,
		amount:         c.amount,
		accountID:      c.accountID,
		idempotenceKey: c.idempotenceKey,
	}
}

func (c CashIn) NewAuthorizationInput() AuthorizationInput {
	return AuthorizationInput{
		operation:      "CREDIT",
		currency:       c.currency,
		amount:         c.amount,
		accountID:      c.accountID,
		idempotenceKey: c.idempotenceKey,
	}
}

func (c CashIn) NewTransactionInput(authorizationID int) TransactionInput {
	return TransactionInput{
		authorizationID: authorizationID,
		operation:       "CREDIT",
		currency:        c.currency,
		amount:          c.amount,
		accountID:       c.accountID,
		descriptor:      "New cash-in",
	}
}

func (c CashIn) NewStreamInput(authorizationID, transactionID int) StreamInput {
	return StreamInput{
		authorizationID: authorizationID,
		transactionID:   transactionID,
		operation:       "CREDIT",
		currency:        c.currency,
		amount:          c.amount,
		accountID:       c.accountID,
		descriptor:      "New cash-in",
	}
}
