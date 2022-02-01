package usecase

import (
	"context"
	"log"

	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"go-opentelemetry-example/domain"
)

type CreateFinancialOperationUseCase interface {
	Execute(context.Context, CreateFinancialOperationInput) (CreateFinancialOperationOutput, error)
}

type Authorizer interface {
	Create(context.Context, AuthorizationInput) error
	FindByIdempotenceKey(context.Context, AuthorizationInput) (domain.Authorization, error)
}

type RulesEvaluator interface {
	SpendingControls(context.Context, RulesInput) error
}

type Ledger interface {
	AvailableEntries(context.Context, LedgerInput) error
}

type Acquirer interface {
	Sale(context.Context, string, string, string, string) error
}

type Transactor interface {
	Transaction(context.Context, TransactionInput) error
}

type PaymentStreamer interface {
	Publish(context.Context)
}

type CreateFinancialOperationInputer interface {
	Build(context.Context) CreateFinancialOperationInput
}

type CreateFinancialOperationInput struct {
	ID             string `json:"id"`
	Amount         int
	External       bool
	IdempotenceKey string `json:"idempotence_key"`
	CashIn         domain.CashIn
	CashOut        domain.CashOut
	Transfer       domain.Transfer
}

type CreateFinancialOperationOutput struct {
	ID             string `json:"id"`
	IdempotenceKey string `json:"idempotence_key"`
}

type createFinancialOperationOrchestrator struct {
	authorizer     Authorizer
	rulesEvaluator RulesEvaluator
	ledger         Ledger
	acquirer       Acquirer
	transactor     Transactor
	stream         PaymentStreamer
	tracer         trace.Tracer
}

func NewCreateFinancialOperationOrchestrator(
	authorizer Authorizer,
	rulesEvaluation RulesEvaluator,
	ledger Ledger,
	acquirer Acquirer,
	transactor Transactor,
	stream PaymentStreamer,
	tracer trace.Tracer,
) CreateFinancialOperationUseCase {
	return createFinancialOperationOrchestrator{
		authorizer:     authorizer,
		rulesEvaluator: rulesEvaluation,
		ledger:         ledger,
		acquirer:       acquirer,
		transactor:     transactor,
		stream:         stream,
		tracer:         tracer,
	}
}

type (
	LedgerInput        struct{}
	RulesInput         struct{}
	AuthorizationInput struct{}
	TransactionInput   struct{}
)

func NewLedgerInput(input CreateFinancialOperationInput) LedgerInput {
	return LedgerInput{}
}

func NewRulesInput(input CreateFinancialOperationInput) RulesInput {
	return RulesInput{}
}

func NewAuthorizationInput(input CreateFinancialOperationInput) AuthorizationInput {
	return AuthorizationInput{}
}

func NewTransactionInput(input CreateFinancialOperationInput) TransactionInput {
	return TransactionInput{}
}

func (c createFinancialOperationOrchestrator) Execute(
	ctx context.Context,
	input CreateFinancialOperationInput,
) (CreateFinancialOperationOutput, error) {
	ctx, span := c.tracer.Start(ctx, "usecase::create_user")
	defer span.End()

	// Start financial operation

	// First step - Authorization
	authorization, err := c.authorizer.FindByIdempotenceKey(ctx, NewAuthorizationInput(input))
	if err != nil {
		return CreateFinancialOperationOutput{}, err
	}

	if authorization.AlreadyExists() {
		return CreateFinancialOperationOutput{
			ID:             authorization.ID(),
			IdempotenceKey: authorization.IdempotenceKey(),
		}, domain.ErrDuplicatedAuthorization
	}

	if err := c.rulesEvaluator.SpendingControls(ctx, NewRulesInput(input)); err != nil {
		return CreateFinancialOperationOutput{}, err
	}

	if err := c.ledger.AvailableEntries(ctx, NewLedgerInput(input)); err != nil {
		return CreateFinancialOperationOutput{}, err
	}

	err = c.authorizer.Create(ctx, NewAuthorizationInput(input))
	if err != nil {
		return CreateFinancialOperationOutput{}, err
	}

	// Second step - Transaction
	err = c.transactor.Transaction(ctx, NewTransactionInput(input))
	if err != nil {
		return CreateFinancialOperationOutput{}, err
	}

	// Data stream step
	c.stream.Publish(ctx)

	log.Print("UseCase execute success")

	span.SetStatus(otelcodes.Ok, "UseCase execute success")

	// Finish financial operation

	return CreateFinancialOperationOutput{}, nil
}
