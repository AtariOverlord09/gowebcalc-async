package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	ID string `json:"id"`
}

type ExpressionResponse struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

type ExpressionsListResponse struct {
	Expressions []ExpressionResponse `json:"expressions"`
}

type TaskResponse struct {
	Task struct {
		ID            string `json:"id"`
		Arg1          string `json:"arg1"`
		Arg2          string `json:"arg2"`
		Operation     string `json:"operation"`
		OperationTime int    `json:"operation_time"`
	} `json:"task"`
}

type TaskResultRequest struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

type CalcController struct {
	calcService CalculatorService
}

func NewCalcController(service CalculatorService) *CalcController {
	return &CalcController{
		calcService: service,
	}
}

func (c *CalcController) SubmitCalculation(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	id, err := c.calcService.SubmitExpression(req.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := CalculateResponse{ID: id}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c *CalcController) GetExpressions(w http.ResponseWriter, r *http.Request) {
	expressions, err := c.calcService.GetAllExpressions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ExpressionsListResponse{Expressions: expressions}
	json.NewEncoder(w).Encode(response)
}

func (c *CalcController) GetExpression(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	expression, err := c.calcService.GetExpression(id)
	if err != nil {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]ExpressionResponse{
		"expression": expression,
	})
}

func (c *CalcController) GetTask(w http.ResponseWriter, r *http.Request) {
	task, err := c.calcService.GetNextTask()
	if err != nil {
		http.Error(w, "No tasks available", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func (c *CalcController) SubmitTaskResult(w http.ResponseWriter, r *http.Request) {
	var req TaskResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	if err := c.calcService.SubmitTaskResult(req.ID, req.Result); err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type CalculatorService interface {
	SubmitExpression(expression string) (string, error)
	GetAllExpressions() ([]ExpressionResponse, error)
	GetExpression(id string) (ExpressionResponse, error)
	GetNextTask() (*TaskResponse, error)
	SubmitTaskResult(id string, result float64) error
}
