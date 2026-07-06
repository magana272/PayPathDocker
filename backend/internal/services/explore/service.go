package explore

import (
	"paypath/internal/services/debts"
	"paypath/internal/services/expenses"
	"paypath/internal/services/income"
	"paypath/internal/services/reporting"
	"paypath/pkg/utils"
)

type Service struct {
	reports *reporting.Repository
}

func NewService(reports *reporting.Repository) *Service {
	return &Service{reports: reports}
}

func (s *Service) Get(uid int) (Data, error) {
	incomes, exps, dbts, liqs, err := s.reports.Financials(uid)
	if err != nil {
		return Data{}, err
	}
	summary := reporting.BuildSummary(incomes, exps, dbts, liqs)
	taxes := income.CalcTaxes(income.CalcAnnualGross(incomes))
	surplus := taxes.MonthlyNet - expenses.CalcMonthly(exps)
	payoff := debts.RunAvalanche(dbts, surplus)
	scenarios, err := s.reports.Scenarios(uid)
	if err != nil {
		return Data{}, err
	}
	annualGross := income.CalcAnnualGross(incomes)
	monthlyExp := expenses.CalcMonthly(exps)
	cashflow := reporting.ProjectCashflow(liqs, incomes, exps, annualGross, monthlyExp, 90)
	return Data{
		Summary:   summary,
		Payoff:    payoff,
		Scenarios: utils.NonNil(scenarios),
		Debts:     utils.NonNil(dbts),
		Cashflow:  cashflow,
	}, nil
}
