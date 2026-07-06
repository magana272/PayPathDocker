package handler

import (
	"net/http"

	"paypath/internal/middleware"
	"paypath/internal/services/ai/strategies"
	"paypath/pkg/response"
)

type StrategiesHandler struct {
	svc *strategies.Service
}

func NewStrategiesHandler(svc *strategies.Service) StrategiesHandler {
	return StrategiesHandler{svc: svc}
}

func (h StrategiesHandler) DebtPayoffStrategy(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.DebtPayoff(r.Context(), middleware.UserID(r))
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	response.JSON(w, 200, out)
}

func (h StrategiesHandler) SavingsPlan(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.SavingsPlan(r.Context(), middleware.UserID(r))
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	response.JSON(w, 200, out)
}

func (h StrategiesHandler) ExpenseAudit(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.ExpenseAudit(r.Context(), middleware.UserID(r))
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	response.JSON(w, 200, out)
}

func (h StrategiesHandler) IncomeBoost(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.IncomeBoost(r.Context(), middleware.UserID(r))
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	response.JSON(w, 200, out)
}
