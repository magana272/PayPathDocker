package handler

import (
	"net/http"

	"paypath/internal/middleware"
	"paypath/internal/services/dashboard"
	"paypath/internal/services/explore"
	"paypath/internal/services/settings"
	"paypath/pkg/response"
)

type BundleHandler struct {
	dashboard *dashboard.Service
	explore   *explore.Service
	settings  *settings.Service
}

func NewBundleHandler(d *dashboard.Service, e *explore.Service, s *settings.Service) BundleHandler {
	return BundleHandler{dashboard: d, explore: e, settings: s}
}

func (h BundleHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	out, err := h.dashboard.Get(middleware.UserID(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response.JSON(w, 200, out)
}

func (h BundleHandler) Explore(w http.ResponseWriter, r *http.Request) {
	out, err := h.explore.Get(middleware.UserID(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response.JSON(w, 200, out)
}

func (h BundleHandler) Settings(w http.ResponseWriter, r *http.Request) {
	out, err := h.settings.Get(middleware.UserID(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response.JSON(w, 200, out)
}
