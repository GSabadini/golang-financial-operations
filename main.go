package main

import (
	"fmt"
	"log"
	"net/http"

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

	//tracer := otel.Tracer(DefaultAppName)

	//createUserHandler := handler.NewCreateUser(
	//	usecase.NewCreateAccount(
	//		repository.NewCreateAccount(tracer),
	//		tracer,
	//	),
	//	tracer,
	//)

	//r.HandleFunc("/users", createUserHandler.Handle).Methods(http.MethodPost)

	log.Print("Start server in port:", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
	if err != nil {
		log.Fatalln("Error start server", err)
		return
	}
}
