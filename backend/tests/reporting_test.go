package tests

import (
	"testing"

	"paypath/internal/liquid"
	"paypath/internal/services/debts"
	"paypath/internal/services/expenses"
	"paypath/internal/services/income"
	"paypath/internal/services/reporting"
)

func TestBuildSummary(t *testing.T) {
	incomes := []income.Income{{PayType: "salary", AnnualSalary: fptr(120000)}}
	exps := []expenses.Expense{{Cost: 2000, Frequency: "monthly"}}
	dbts := []debts.Debt{{Name: "cc", Balance: 10000, APY: 20, Type: "credit_card"}}
	liqs := []liquid.Liquid{{Balance: 5000}}

	s := reporting.BuildSummary(incomes, exps, dbts, liqs)
	if !approx(s.MonthlyGross, 10000, 0.01) {
		t.Errorf("MonthlyGross = %v, want 10000", s.MonthlyGross)
	}
	if !approx(s.TotalDebt, 10000, 0.01) {
		t.Errorf("TotalDebt = %v, want 10000", s.TotalDebt)
	}
	if !approx(s.TotalLiquid, 5000, 0.01) {
		t.Errorf("TotalLiquid = %v, want 5000", s.TotalLiquid)
	}
	if !approx(s.NetWorth, -5000, 0.01) {
		t.Errorf("NetWorth = %v, want -5000", s.NetWorth)
	}
	if !approx(s.MonthlyExpenses, 2000, 0.01) {
		t.Errorf("MonthlyExpenses = %v, want 2000", s.MonthlyExpenses)
	}
}

func TestSolveScenarios(t *testing.T) {
	dbts := []debts.Debt{{Name: "cc", Balance: 5000, APY: 20, Type: "credit_card"}}
	exps := []expenses.Expense{{Cost: 1000, Frequency: "monthly"}}
	incomes := []income.Income{{PayType: "hourly", PayPerHour: fptr(20), HourPerDay: fptr(8)}}

	scs := reporting.SolveScenarios(dbts, exps, incomes)
	if len(scs) != 7 {
		t.Fatalf("want 7 scenarios, got %d", len(scs))
	}
	for _, sc := range scs {
		if sc.HourlyRate < 0 {
			t.Errorf("negative hourly rate for %d-month scenario", sc.Months)
		}
	}
}

func TestBuildCalendarPaydays(t *testing.T) {
	incomes := []income.Income{{Job: "J", PayType: "salary", AnnualSalary: fptr(120000)}}
	cal := reporting.BuildCalendar(2025, 6, nil, incomes, income.CalcAnnualGross(incomes))
	if cal.Year != 2025 || cal.Month != 6 {
		t.Fatalf("year/month = %d/%d", cal.Year, cal.Month)
	}
	if len(cal.Events["1"]) == 0 || len(cal.Events["15"]) == 0 {
		t.Fatalf("expected semi-monthly paydays on days 1 and 15")
	}
}

func TestBuildCalendarDueDateClamp(t *testing.T) {
	due := 31
	exps := []expenses.Expense{{Expense: "rent", Cost: 1000, Frequency: "monthly", DueDate: &due}}
	cal := reporting.BuildCalendar(2025, 2, exps, nil, 0)
	if len(cal.Events["28"]) == 0 {
		t.Fatalf("due_date 31 in February should clamp to day 28")
	}
}

func TestProjectCashflowLength(t *testing.T) {
	cf := reporting.ProjectCashflow([]liquid.Liquid{{Balance: 1000}}, nil, nil, 0, 0, 30)
	if len(cf) != 30 {
		t.Fatalf("want 30 days, got %d", len(cf))
	}
}
