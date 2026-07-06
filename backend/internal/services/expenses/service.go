package expenses

import "paypath/internal/storage/cache"

func CalcMonthly(expenses []Expense) float64 {
	total := 0.0
	for _, e := range expenses {
		total += Normalize(e)
	}
	return total
}

func Normalize(e Expense) float64 {
	switch e.Frequency {
	case "one-time":
		return 0
	case "biweekly":
		return e.Cost * 26.0 / 12.0
	case "weekly":
		return e.Cost * 52.0 / 12.0
	case "yearly":
		return e.Cost / 12.0
	default:
		return e.Cost
	}
}

type Service struct {
	repo    Repository
	derived *cache.DerivedCache
}

func NewService(repo Repository, derived *cache.DerivedCache) *Service {
	return &Service{repo: repo, derived: derived}
}

func (s *Service) List(userID int) ([]Expense, error) {
	return s.repo.All(userID)
}

func (s *Service) Create(userID int, e Expense) (Expense, error) {
	out, err := s.repo.Create(userID, e)
	if err == nil {
		s.derived.Invalidate(userID)
	}
	return out, err
}

func (s *Service) Update(userID, id int, e Expense) (*Expense, error) {
	out, err := s.repo.Update(userID, id, e)
	if err == nil && out != nil {
		s.derived.Invalidate(userID)
	}
	return out, err
}

func (s *Service) Delete(userID, id int) (bool, error) {
	ok, err := s.repo.Delete(userID, id)
	if err == nil && ok {
		s.derived.Invalidate(userID)
	}
	return ok, err
}
