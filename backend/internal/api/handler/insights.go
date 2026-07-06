package handler

import (
	"errors"
	"net/http"

	"paypath/internal/clients"
	"paypath/internal/middleware"
	"paypath/internal/services/ai/insights"
	"paypath/pkg/response"
)

type InsightsHandler struct {
	svc *insights.Service
}

func NewInsightsHandler(svc *insights.Service) InsightsHandler {
	return InsightsHandler{svc: svc}
}

func (h InsightsHandler) GetInsights(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.Get(r.Context(), middleware.UserID(r))
	if err != nil {
		if errors.Is(err, clients.ErrNoAPIKey) {
			http.Error(w, "OPENAI_API_KEY not set", 500)
			return
		}
		http.Error(w, err.Error(), 502)
		return
	}
	response.JSON(w, 200, out)
}
