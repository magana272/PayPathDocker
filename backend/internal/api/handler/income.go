package handler

import (
	"net/http"
	"strconv"

	"paypath/internal/middleware"
	"paypath/internal/services/income"
	"paypath/pkg/response"
)

type IncomeHandler struct {
	svc *income.Service
}

func NewIncomeHandler(svc *income.Service) IncomeHandler {
	return IncomeHandler{svc: svc}
}

func (h IncomeHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.List(middleware.UserID(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if list == nil {
		list = []income.Income{}
	}
	response.JSON(w, 200, list)
}

func (h IncomeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var inc income.Income
	if err := response.Decode(r, &inc); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	created, err := h.svc.Create(middleware.UserID(r), inc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response.JSON(w, 201, created)
}

func (h IncomeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	var inc income.Income
	if err := response.Decode(r, &inc); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	inc.ID = id
	updated, err := h.svc.Update(middleware.UserID(r), id, inc)
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

func (h IncomeHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
