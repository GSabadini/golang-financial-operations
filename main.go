package main

import (
	"fmt"
	"go-opentelemetry-example/adapter"
	"go-opentelemetry-example/handler"
	"go-opentelemetry-example/infrastructure"
	"go-opentelemetry-example/usecase"
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go-opentelemetry-example/infrastructure/opentelemetry"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

const (
	PORT           = 8080
	DefaultAppName = ""
)

func main() {
	shutdownTracer := opentelemetry.NewTracer()
	defer shutdownTracer()

	r := mux.NewRouter()
	r.Use(otelmux.Middleware(DefaultAppName))

	tracer := otel.Tracer(DefaultAppName)

	httpClient := infrastructure.NewClient(infrastructure.NewRequest())

	createFinancialOperationHandler := handler.NewCreateFinancialOperation(
		usecase.NewCreateFinancialOperationOrchestrator(
			adapter.NewAuthorization(httpClient, os.Getenv("AUTHORIZATION_URI")),
			adapter.NewRulesEvaluator(httpClient, os.Getenv("RULES_EVALUATOR_URI")),
			adapter.NewLedger(httpClient, os.Getenv("LEDGER_URI")),
			adapter.NewTransaction(httpClient, os.Getenv("TRANSACTION_URI")),
			adapter.NewStream(),
			tracer,
		),
		tracer,
	)

	r.HandleFunc("/v1/financial-operations", createFinancialOperationHandler.Handle).Methods(http.MethodPost)

	log.Print("Start server in port:", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
	if err != nil {
		log.Fatalln("Error start server", err)
		return
	}
}
