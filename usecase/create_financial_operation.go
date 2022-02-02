package usecase

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/trace"

	"go-opentelemetry-example/domain"
)

type CreateFinancialOperationUseCase interface {
	Execute(context.Context, domain.FinancialOperator) (CreateFinancialOperationOutput, error)
}

type CreateFinancialOperationInput struct {
	CashIn   domain.CashIn
	CashOut  domain.CashOut
	Transfer domain.Transfer
}

type CreateFinancialOperationOutput struct {
	ID             string `json:"id"`
	IdempotenceKey string `json:"idempotence_key"`
}

type createFinancialOperationOrchestrator struct {
	authorizer     domain.Authorizer
	rulesEvaluator domain.RulesEvaluator
	ledger         domain.Ledger
	acquirer       domain.Acquirer
	transactor     domain.Transactor
	streamer       domain.FinancialOperationStreamer
	tracer         trace.Tracer
}

func NewCreateFinancialOperationOrchestrator(
	authorizer domain.Authorizer,
	rulesEvaluation domain.RulesEvaluator,
	ledger domain.Ledger,
	//acquirer domain.Acquirer,
	transactor domain.Transactor,
	streamer domain.FinancialOperationStreamer,
	tracer trace.Tracer,
) CreateFinancialOperationUseCase {
	return createFinancialOperationOrchestrator{
		authorizer:     authorizer,
		rulesEvaluator: rulesEvaluation,
		ledger:         ledger,
		//acquirer:       acquirer,
		transactor: transactor,
		streamer:   streamer,
		tracer:     tracer,
	}
}

func (c createFinancialOperationOrchestrator) Execute(
	ctx context.Context,
	financialOperator domain.FinancialOperator,
) (CreateFinancialOperationOutput, error) {
	ctx, span := c.tracer.Start(ctx, "usecase::create_financial_operator")
	defer span.End()

	log.Print("Start financial operation")

	financialOperation, err := financialOperator.Process(
		ctx,
		c.authorizer,
		c.rulesEvaluator,
		c.ledger,
		c.transactor,
		c.streamer,
	)
	if err != nil {
		log.Print("Error financial operation", err)

		span.RecordError(err)

		return CreateFinancialOperationOutput{}, err
	}

	log.Print("Finish financial operation")

	return CreateFinancialOperationOutput{
		ID:             financialOperation.ID(),
		IdempotenceKey: financialOperation.IdempotenceKey(),
	}, nil
}
