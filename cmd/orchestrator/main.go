package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "path/to/your/rest"
)

func main() {
    router := mux.NewRouter()
    
    // Initialize your service implementation
    calcService := NewCalculatorService() // Implement this based on your domain logic
    controller := rest.NewCalcController(calcService)

    // Public API routes
    api := router.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/calculate", controller.SubmitCalculation).Methods(http.MethodPost)
    api.HandleFunc("/expressions", controller.GetExpressions).Methods(http.MethodGet)
    api.HandleFunc("/expressions/{id}", controller.GetExpression).Methods(http.MethodGet)

    // Internal API routes
    internal := router.PathPrefix("/internal").Subrouter()
    internal.HandleFunc("/task", controller.GetTask).Methods(http.MethodGet)
    internal.HandleFunc("/task", controller.SubmitTaskResult).Methods(http.MethodPost)

    log.Fatal(http.ListenAndServe(":8080", router))
}
