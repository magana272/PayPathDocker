package settings

import (
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
	return Data{
		Income:   utils.NonNil(incomes),
		Liquid:   utils.NonNil(liqs),
		Expenses: utils.NonNil(exps),
		Debts:    utils.NonNil(dbts),
	}, nil
}
