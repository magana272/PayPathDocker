package debts

import (
	"math"
	"sort"

	"paypath/internal/storage/cache"
	"paypath/pkg/utils"
)

type PayoffResult struct {
	Budget        float64                  `json:"budget"`
	Months        int                      `json:"months"`
	TotalInterest float64                  `json:"total_interest,omitempty"`
	History       []map[string]interface{} `json:"history"`
}

func minPayment(balance, apy float64, debtType string) float64 {
	monthlyRate := apy / 100.0 / 12.0
	interest := balance * monthlyRate

	switch debtType {
	case "credit_card":
		return math.Max(balance*0.01+interest, 25)
	default:
		if monthlyRate == 0 {
			return balance / 60
		}
		n := 120.0
		if debtType == "car" {
			n = 60
		}
		return balance * monthlyRate * math.Pow(1+monthlyRate, n) / (math.Pow(1+monthlyRate, n) - 1)
	}
}

func TotalMinPayments(debts []Debt) float64 {
	total := 0.0
	for _, d := range debts {
		total += minPayment(d.Balance, d.APY, d.Type)
	}
	return total
}

type debtState struct {
	name     string
	balance  float64
	apy      float64
	debtType string
	fixedMin float64
	min      float64
}

func RunAvalanche(debts []Debt, budget float64) PayoffResult {
	if len(debts) == 0 || budget <= 0 {
		return PayoffResult{Budget: utils.Round2(budget), Months: 0, History: []map[string]interface{}{}}
	}

	states := make([]debtState, 0, len(debts))
	for _, d := range debts {
		if d.Balance > 0 {
			fm := 0.0
			if d.Type != "credit_card" {
				fm = minPayment(d.Balance, d.APY, d.Type)
			}
			states = append(states, debtState{
				name: d.Name, balance: d.Balance, apy: d.APY, debtType: d.Type, fixedMin: fm,
			})
		}
	}

	if len(states) == 0 {
		return PayoffResult{Budget: utils.Round2(budget), Months: 0, History: []map[string]interface{}{}}
	}

	sort.Slice(states, func(i, j int) bool {
		return states[i].apy > states[j].apy
	})

	var history []map[string]interface{}
	cumulativeInterest := 0.0

	for month := 1; month <= 480; month++ {
		sumMins := 0.0
		for i := range states {
			if states[i].balance <= 0 {
				states[i].min = 0
				continue
			}
			if states[i].debtType == "credit_card" {
				states[i].min = minPayment(states[i].balance, states[i].apy, states[i].debtType)
			} else {
				states[i].min = states[i].fixedMin
			}
			interest := states[i].balance * states[i].apy / 100.0 / 12.0
			cumulativeInterest += interest
			states[i].balance += interest
			states[i].min = math.Min(states[i].min, states[i].balance)
			sumMins += states[i].min
		}

		remaining := math.Max(budget, sumMins)
		for i := range states {
			if states[i].balance <= 0 {
				continue
			}
			payment := math.Min(states[i].min, remaining)
			states[i].balance -= payment
			remaining -= payment
			if states[i].balance < 0.01 {
				states[i].balance = 0
			}
		}

		for i := range states {
			if remaining <= 0 {
				break
			}
			if states[i].balance <= 0 {
				continue
			}
			payment := math.Min(remaining, states[i].balance)
			states[i].balance -= payment
			remaining -= payment
			if states[i].balance < 0.01 {
				states[i].balance = 0
			}
		}

		entry := map[string]interface{}{"month": month}
		total := 0.0
		for _, s := range states {
			entry[s.name] = utils.Round2(s.balance)
			total += s.balance
		}
		entry["total"] = utils.Round2(total)
		entry["interest"] = utils.Round2(cumulativeInterest)
		history = append(history, entry)

		if total <= 0 {
			return PayoffResult{Budget: utils.Round2(budget), Months: month, TotalInterest: utils.Round2(cumulativeInterest), History: history}
		}
	}

	return PayoffResult{Budget: utils.Round2(budget), Months: 480, TotalInterest: utils.Round2(cumulativeInterest), History: history}
}

func RunSnowball(debts []Debt, budget float64) PayoffResult {
	if len(debts) == 0 || budget <= 0 {
		return PayoffResult{Budget: utils.Round2(budget), Months: 0, History: []map[string]interface{}{}}
	}

	states := make([]debtState, 0, len(debts))
	for _, d := range debts {
		if d.Balance > 0 {
			fm := 0.0
			if d.Type != "credit_card" {
				fm = minPayment(d.Balance, d.APY, d.Type)
			}
			states = append(states, debtState{
				name: d.Name, balance: d.Balance, apy: d.APY, debtType: d.Type, fixedMin: fm,
			})
		}
	}

	if len(states) == 0 {
		return PayoffResult{Budget: utils.Round2(budget), Months: 0, History: []map[string]interface{}{}}
	}

	sort.Slice(states, func(i, j int) bool {
		return states[i].balance < states[j].balance
	})

	var history []map[string]interface{}
	totalInterest := 0.0

	for month := 1; month <= 480; month++ {
		sumMins := 0.0
		for i := range states {
			if states[i].balance <= 0 {
				states[i].min = 0
				continue
			}
			if states[i].debtType == "credit_card" {
				states[i].min = minPayment(states[i].balance, states[i].apy, states[i].debtType)
			} else {
				states[i].min = states[i].fixedMin
			}
			interest := states[i].balance * states[i].apy / 100.0 / 12.0
			totalInterest += interest
			states[i].balance += interest
			states[i].min = math.Min(states[i].min, states[i].balance)
			sumMins += states[i].min
		}

		remaining := math.Max(budget, sumMins)
		for i := range states {
			if states[i].balance <= 0 {
				continue
			}
			payment := math.Min(states[i].min, remaining)
			states[i].balance -= payment
			remaining -= payment
			if states[i].balance < 0.01 {
				states[i].balance = 0
			}
		}

		for i := range states {
			if remaining <= 0 {
				break
			}
			if states[i].balance <= 0 {
				continue
			}
			payment := math.Min(remaining, states[i].balance)
			states[i].balance -= payment
			remaining -= payment
			if states[i].balance < 0.01 {
				states[i].balance = 0
			}
		}

		entry := map[string]interface{}{"month": month}
		total := 0.0
		for _, s := range states {
			entry[s.name] = utils.Round2(s.balance)
			total += s.balance
		}
		entry["total"] = utils.Round2(total)
		entry["interest"] = utils.Round2(totalInterest)
		history = append(history, entry)

		if total <= 0 {
			return PayoffResult{Budget: utils.Round2(budget), Months: month, History: history, TotalInterest: utils.Round2(totalInterest)}
		}
	}

	return PayoffResult{Budget: utils.Round2(budget), Months: 480, History: history, TotalInterest: utils.Round2(totalInterest)}
}

func RunAvalancheWithInterest(debts []Debt, budget float64) PayoffResult {
	return RunAvalanche(debts, budget)
}

type Service struct {
	repo    Repository
	derived *cache.DerivedCache
}

func NewService(repo Repository, derived *cache.DerivedCache) *Service {
	return &Service{repo: repo, derived: derived}
}

func (s *Service) List(userID int) ([]Debt, error) {
	return s.repo.All(userID)
}

func (s *Service) Create(userID int, d Debt) (Debt, error) {
	out, err := s.repo.Create(userID, d)
	if err == nil {
		s.derived.Invalidate(userID)
	}
	return out, err
}

func (s *Service) Update(userID, id int, d Debt) (bool, error) {
	ok, err := s.repo.Update(userID, id, d)
	if err == nil && ok {
		s.derived.Invalidate(userID)
	}
	return ok, err
}

func (s *Service) Delete(userID, id int) (bool, error) {
	ok, err := s.repo.Delete(userID, id)
	if err == nil && ok {
		s.derived.Invalidate(userID)
	}
	return ok, err
}
