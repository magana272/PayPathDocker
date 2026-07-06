package handler

import (
	"net/http"
	"strconv"

	"paypath/internal/middleware"
	"paypath/internal/services/expenses"
	"paypath/pkg/response"
)

type ExpensesHandler struct {
	svc *expenses.Service
}

func NewExpensesHandler(svc *expenses.Service) ExpensesHandler {
	return ExpensesHandler{svc: svc}
}

func (h ExpensesHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.List(middleware.UserID(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if list == nil {
		list = []expenses.Expense{}
	}
	response.JSON(w, 200, list)
}

func (h ExpensesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var e expenses.Expense
	if err := response.Decode(r, &e); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	created, err := h.svc.Create(middleware.UserID(r), e)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response.JSON(w, 201, created)
}

func (h ExpensesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	var e expenses.Expense
	if err := response.Decode(r, &e); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	e.ID = id
	updated, err := h.svc.Update(middleware.UserID(r), id, e)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if updated == nil {
		http.Error(w, "not found", 404)
		return
	}
	response.JSON(w, 200, updated)
}

func (h ExpensesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	found, err := h.svc.Delete(middleware.UserID(r), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if !found {
		http.Error(w, "not found", 404)
		return
	}
	response.JSON(w, 200, map[string]interface{}{})
}
