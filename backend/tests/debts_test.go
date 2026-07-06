package tests

import (
	"testing"

	"paypath/internal/services/debts"
)

func TestRunAvalancheNoDebts(t *testing.T) {
	if r := debts.RunAvalanche(nil, 500); r.Months != 0 {
		t.Fatalf("no debts should give 0 months, got %d", r.Months)
	}
}

func TestRunAvalancheZeroBudget(t *testing.T) {
	ds := []debts.Debt{{Name: "cc", Balance: 1000, APY: 20, Type: "credit_card"}}
	if r := debts.RunAvalanche(ds, 0); r.Months != 0 {
		t.Fatalf("zero budget should give 0 months, got %d", r.Months)
	}
}

func TestRunAvalanchePaysOff(t *testing.T) {
	ds := []debts.Debt{{Name: "cc", Balance: 1000, APY: 20, Type: "credit_card"}}
	r := debts.RunAvalanche(ds, 500)
	if r.Months <= 0 {
		t.Fatalf("should pay off, got %d months", r.Months)
	}
	last := r.History[len(r.History)-1]
	if total, _ := last["total"].(float64); total > 0.01 {
		t.Fatalf("final total should be ~0, got %v", total)
	}
}

func TestRunAvalancheClearsHighAPYFirst(t *testing.T) {
	ds := []debts.Debt{
		{Name: "low", Balance: 1000, APY: 5, Type: "credit_card"},
		{Name: "high", Balance: 1000, APY: 25, Type: "credit_card"},
	}
	r := debts.RunAvalanche(ds, 600)
	high, low := monthClearedTo(r, "high"), monthClearedTo(r, "low")
	if high == 0 || low == 0 || high > low {
		t.Fatalf("avalanche should clear high-APY first: high@%d low@%d", high, low)
	}
}

func TestRunSnowballClearsSmallestFirst(t *testing.T) {
	ds := []debts.Debt{
		{Name: "big", Balance: 5000, APY: 25, Type: "credit_card"},
		{Name: "small", Balance: 500, APY: 5, Type: "credit_card"},
	}
	r := debts.RunSnowball(ds, 700)
	small, big := monthClearedTo(r, "small"), monthClearedTo(r, "big")
	if small == 0 || big == 0 || small > big {
		t.Fatalf("snowball should clear smallest balance first: small@%d big@%d", small, big)
	}
}

func monthClearedTo(r debts.PayoffResult, name string) int {
	for i, e := range r.History {
		if v, ok := e[name].(float64); ok && v == 0 {
			return i + 1
		}
	}
	return 0
}
