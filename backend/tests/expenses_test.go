package tests

import (
	"testing"

	"paypath/internal/services/expenses"
)

func TestNormalize(t *testing.T) {
	cases := []struct {
		freq string
		cost float64
		want float64
	}{
		{"one-time", 100, 0},
		{"weekly", 100, 100 * 52.0 / 12.0},
		{"biweekly", 100, 100 * 26.0 / 12.0},
		{"yearly", 1200, 100},
		{"monthly", 100, 100},
		{"", 100, 100},
	}
	for _, c := range cases {
		got := expenses.Normalize(expenses.Expense{Cost: c.cost, Frequency: c.freq})
		if !approx(got, c.want, 0.001) {
			t.Errorf("Normalize(%q, %v) = %v, want %v", c.freq, c.cost, got, c.want)
		}
	}
}

func TestCalcMonthly(t *testing.T) {
	exps := []expenses.Expense{
		{Cost: 100, Frequency: "monthly"},
		{Cost: 1200, Frequency: "yearly"},
		{Cost: 500, Frequency: "one-time"},
	}
	if got := expenses.CalcMonthly(exps); !approx(got, 200, 0.001) {
		t.Fatalf("CalcMonthly = %v, want 200", got)
	}
}
