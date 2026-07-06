package handler

import (
	"net/http"
	"strconv"

	"paypath/internal/liquid"
	"paypath/internal/middleware"
	"paypath/pkg/response"
)

type LiquidHandler struct {
	repo liquid.Repository
}

func NewLiquidHandler(repo liquid.Repository) LiquidHandler {
	return LiquidHandler{repo: repo}
}

func (h LiquidHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.All(middleware.UserID(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if list == nil {
		list = []liquid.Liquid{}
	}
	response.JSON(w, 200, list)
}

func (h LiquidHandler) Create(w http.ResponseWriter, r *http.Request) {
	var l liquid.Liquid
	if err := response.Decode(r, &l); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	created, err := h.repo.Create(middleware.UserID(r), l)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response.JSON(w, 201, created)
}

func (h LiquidHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	var l liquid.Liquid
	if err := response.Decode(r, &l); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	l.ID = id
	found, err := h.repo.Update(middleware.UserID(r), id, l)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if !found {
		http.Error(w, "not found", 404)
		return
	}
	response.JSON(w, 200, l)
}
