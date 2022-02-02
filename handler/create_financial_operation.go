package handler

import (
	"encoding/json"
	"errors"
	"go-opentelemetry-example/domain"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"

	"go-opentelemetry-example/usecase"
)

type CreateFinancialOperation struct {
	uc     usecase.CreateFinancialOperationUseCase
	tracer trace.Tracer
}

func NewCreateFinancialOperation(uc usecase.CreateFinancialOperationUseCase, tracer trace.Tracer) CreateFinancialOperation {
	return CreateFinancialOperation{
		uc:     uc,
		tracer: tracer,
	}
}

func (c CreateFinancialOperation) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, span := c.tracer.Start(r.Context(), "handler::create_financial_operation")
	defer span.End()

	var genericInput GenericFinancialOperationInput
	if err := json.NewDecoder(r.Body).Decode(&genericInput); err != nil {
		log.Print("Handler execute error", err)

		span.RecordError(err)

		_ = response(w, err, http.StatusBadRequest)
	}
	defer r.Body.Close()

	financialOperator, err := c.factoryFinancialOperator(genericInput)
	if err != nil {
		log.Print("Handler execute error", err)

		span.RecordError(err)

		_ = response(w, err, http.StatusBadRequest)
	}

	output, err := c.uc.Execute(ctx, financialOperator)
	if err != nil {
		log.Print("Handler execute error", err)

		span.RecordError(err)

		_ = response(w, err, http.StatusInternalServerError)
	}

	log.Print("Handler execute success")

	_ = response(w, output, http.StatusCreated)
}

type GenericFinancialOperationInput struct {
	to             To
	from           From
	idempotencyKey string
	descriptor     string
}

func (g GenericFinancialOperationInput) isCashIn() bool {
	return g.to.accountID != 0 && g.from.accountID == 0
}

func (g GenericFinancialOperationInput) isCashOut() bool {
	return g.to.accountID == 0 && g.from.accountID != 0
}

func (g GenericFinancialOperationInput) isTransfer() bool {
	return g.to.accountID != 0 && g.from.accountID != 0
}

type To struct {
	accountID int
	amount    int
	currency  string
}

type From struct {
	accountID int
	amount    int
	currency  string
}

func (c CreateFinancialOperation) factoryFinancialOperator(
	input GenericFinancialOperationInput,
) (domain.FinancialOperator, error) {
	switch {
	case input.isCashIn():
		return domain.NewCashIn(
			input.idempotencyKey,
			input.to.accountID,
			input.to.amount,
			input.to.currency,
		), nil

	case input.isCashOut():
	//	return domain.NewCashOut(
	//		input.idempotencyKey,
	//		input.from.accountID,
	//		input.from.amount,
	//		input.from.currency,
	//	), nil

	case input.isTransfer():
		//	return domain.NewTransfer(
		//		input.idempotencyKey,
		//		input.to.accountID,
		//		input.from.accountID,
		//		input.to.amount,
		//		input.to.currency,
		//	), nil
	}

	return nil, errors.New("not allowed")
}

func response(w http.ResponseWriter, output interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(output)
}
