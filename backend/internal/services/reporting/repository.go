package reporting

import (
	"sync"

	"paypath/internal/liquid"
	"paypath/internal/services/debts"
	"paypath/internal/services/expenses"
	"paypath/internal/services/income"
	"paypath/internal/storage/cache"
)

type Repository struct {
	income   income.Repository
	expenses expenses.Repository
	debts    debts.Repository
	liquid   liquid.Repository
	derived  *cache.DerivedCache
}

func NewRepository(inc income.Repository, exp expenses.Repository, dbt debts.Repository, liq liquid.Repository, derived *cache.DerivedCache) *Repository {
	return &Repository{income: inc, expenses: exp, debts: dbt, liquid: liq, derived: derived}
}

func (r *Repository) Financials(uid int) ([]income.Income, []expenses.Expense, []debts.Debt, []liquid.Liquid, error) {
	var (
		incomes []income.Income
		exps    []expenses.Expense
		dbts    []debts.Debt
		liqs    []liquid.Liquid
		errs    [4]error
		wg      sync.WaitGroup
	)
	wg.Add(4)
	go func() { defer wg.Done(); incomes, errs[0] = r.income.All(uid) }()
	go func() { defer wg.Done(); exps, errs[1] = r.expenses.All(uid) }()
	go func() { defer wg.Done(); dbts, errs[2] = r.debts.All(uid) }()
	go func() { defer wg.Done(); liqs, errs[3] = r.liquid.All(uid) }()
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	return incomes, exps, dbts, liqs, nil
}

func (r *Repository) Scenarios(uid int) ([]Scenario, error) {
	if v, ok := r.derived.Get(uid, "scenarios"); ok {
		return v.([]Scenario), nil
	}
	incomes, exps, dbts, _, err := r.Financials(uid)
	if err != nil {
		return nil, err
	}
	scenarios := SolveScenarios(dbts, exps, incomes)
	r.derived.Set(uid, "scenarios", scenarios)
	return scenarios, nil
}
