package settings

import (
	"paypath/internal/liquid"
	"paypath/internal/services/debts"
	"paypath/internal/services/expenses"
	"paypath/internal/services/income"
)

type Data struct {
	Income   []income.Income    `json:"income"`
	Liquid   []liquid.Liquid    `json:"liquid"`
	Expenses []expenses.Expense `json:"expenses"`
	Debts    []debts.Debt       `json:"debts"`
}
