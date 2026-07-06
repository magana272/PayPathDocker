package handler

import (
	"net/http"
	"strconv"

	"paypath/internal/middleware"
	"paypath/internal/services/debts"
	"paypath/pkg/response"
)

type DebtsHandler struct {
	svc *debts.Service
}

func NewDebtsHandler(svc *debts.Service) DebtsHandler {
	return DebtsHandler{svc: svc}
}

func (h DebtsHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.List(middleware.UserID(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if list == nil {
		list = []debts.Debt{}
	}
	response.JSON(w, 200, list)
}

func (h DebtsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var d debts.Debt
	if err := response.Decode(r, &d); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	created, err := h.svc.Create(middleware.UserID(r), d)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response.JSON(w, 201, created)
}

func (h DebtsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	var d debts.Debt
	if err := response.Decode(r, &d); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	d.ID = id
	found, err := h.svc.Update(middleware.UserID(r), id, d)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if !found {
		http.Error(w, "not found", 404)
		return
	}
	response.JSON(w, 200, d)
}

func (h DebtsHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
