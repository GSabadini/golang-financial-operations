package handler

import (
	"encoding/json"
	"errors"
	"go-opentelemetry-example/domain"
	"log"
	"net/http"
	"time"

	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

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

type GenericFinancialOperationInput struct {
	to             To
	from           From
	idempotencyKey string
	descriptor     string
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

func (c CreateFinancialOperation) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, span := c.tracer.Start(r.Context(), "handler::create_payment")
	defer span.End()

	var genericInput GenericFinancialOperationInput
	if err := json.NewDecoder(r.Body).Decode(&genericInput); err != nil {
		log.Print("Handler execute error", err)

		span.SetStatus(otelcodes.Error, "Handler execute error")
		span.RecordError(err)

		_ = response(w, err, http.StatusBadRequest)
	}
	defer r.Body.Close()

	financialOperation, err := c.strategyBuildFinancialOperation(genericInput)
	if err != nil {
		log.Print("Handler execute error", err)

		span.SetStatus(otelcodes.Error, "Handler execute error")
		span.RecordError(err)

		_ = response(w, err, http.StatusBadRequest)
	}

	output, err := c.uc.Execute(ctx, financialOperation)
	if err != nil {
		log.Print("Handler execute error", err)

		span.SetStatus(otelcodes.Error, "Handler execute error")
		span.RecordError(err)

		_ = response(w, err, http.StatusInternalServerError)
	}

	log.Print("Handler execute success")

	span.SetStatus(otelcodes.Ok, "Handler execute success")

	_ = response(w, output, http.StatusCreated)
}

func (c CreateFinancialOperation) strategyBuildFinancialOperation(
	input GenericFinancialOperationInput,
) (usecase.CreateFinancialOperationInput, error) {
	if input.to != nil && input.from == nil {
		return domain.NewCashIn(
			"",
			input.to.accountID,
			input.idempotencyKey,
			input.to.amount,
			input.to.currency,
			time.Now(),
		), nil
	}

	if input.from != nil && input.to == nil {
		return domain.NewCashOut(
			"",
			input.idempotencyKey,
			input.to.accountID,
			input.to.amount,
			input.to.currency,
			time.Now(),
		), nil
	}

	if input.from != nil && input.to != nil {
		return domain.NewFinancialOperation(
			"",
			input.idempotencyKey,
			input.to.accountID,
			input.from.accountID,
			input.to.amount,
			input.to.currency,
			time.Now(),
		), nil
	}

	return errors.New("not allowed")
}

func response(w http.ResponseWriter, output interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(output)
}
