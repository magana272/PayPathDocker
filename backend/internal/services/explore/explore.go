package explore

import (
	"paypath/internal/services/debts"
	"paypath/internal/services/reporting"
)

type Data struct {
	Summary   reporting.Summary       `json:"summary"`
	Payoff    debts.PayoffResult      `json:"payoff"`
	Scenarios []reporting.Scenario    `json:"scenarios"`
	Debts     []debts.Debt            `json:"debts"`
	Cashflow  []reporting.CashflowDay `json:"cashflow"`
}
