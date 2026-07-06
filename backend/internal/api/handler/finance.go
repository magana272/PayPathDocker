package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"paypath/internal/middleware"
	"paypath/internal/services/reporting"
	"paypath/pkg/response"
)

type FinanceHandler struct {
	svc *reporting.Service
}

func NewFinanceHandler(svc *reporting.Service) FinanceHandler {
	return FinanceHandler{svc: svc}
}

func (h FinanceHandler) Summary(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.Summary(middleware.UserID(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response.JSON(w, 200, out)
}

func (h FinanceHandler) Payoff(w http.ResponseWriter, r *http.Request) {
	extra := 0.0
	if ep := r.URL.Query().Get("extra_payment"); ep != "" {
		if n, err := strconv.ParseFloat(ep, 64); err == nil && n > 0 {
			extra = n
		}
	}
	out, err := h.svc.Payoff(middleware.UserID(r), extra)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response.JSON(w, 200, out)
}

func (h FinanceHandler) Scenarios(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.Scenarios(middleware.UserID(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response.JSON(w, 200, out)
}

func (h FinanceHandler) Cashflow(w http.ResponseWriter, r *http.Request) {
	days := 90
	if d := r.URL.Query().Get("days"); d != "" {
		if n, err := strconv.Atoi(d); err == nil && n > 0 {
			days = n
		}
	}
	out, err := h.svc.Cashflow(middleware.UserID(r), days)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response.JSON(w, 200, out)
}

func (h FinanceHandler) Calendar(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	if y := r.URL.Query().Get("year"); y != "" {
		if n, err := strconv.Atoi(y); err == nil {
			year = n
		}
	}
	if m := r.URL.Query().Get("month"); m != "" {
		if n, err := strconv.Atoi(m); err == nil && n >= 1 && n <= 12 {
			month = n
		}
	}
	cal, err := h.svc.Calendar(middleware.UserID(r), year, month)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	formatted := make(map[string][]reporting.CalendarEvent)
	for day, events := range cal.Events {
		dayNum, _ := strconv.Atoi(day)
		dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, dayNum)
		formatted[dateStr] = events
	}
	response.JSON(w, 200, formatted)
}
