package seed

import (
	"context"

	"paypath/internal/liquid"
	"paypath/internal/services/ai/insights"
	"paypath/internal/services/auth"
	"paypath/internal/services/debts"
	"paypath/internal/services/expenses"
	"paypath/internal/services/income"
	"paypath/internal/storage"
	"paypath/pkg/logger"
	"paypath/pkg/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func Run(db *storage.DB, users auth.Repository, inc income.Repository, exp expenses.Repository, dbt debts.Repository, liq liquid.Repository, ins insights.Repository) {
	count, _ := db.Collection("users").CountDocuments(context.Background(), bson.M{})
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("userpassword"), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error().Err(err).Msg("seed user hash error")
		return
	}
	userID, err := users.CreateUser("user@email.com", string(hash), "Default User")
	if err != nil {
		logger.Log.Error().Err(err).Msg("seed user insert error")
		return
	}
	logger.Log.Info().Msg("seeded default user (user@email.com)")

	seedExpenses(exp, int(userID), "seed/fin_expense.csv")
	seedDebts(dbt, int(userID), "seed/fin_debt.csv")
	seedIncome(inc, int(userID), "seed/fin_income.csv")
	seedLiquid(liq, int(userID), "seed/fin_liqud.csv")
	seedInsights(ins, int(userID))
}

func seedExpenses(repo expenses.Repository, userID int, file string) {
	headers, rows, ok := utils.ReadCSV(file)
	if !ok {
		return
	}
	col := utils.HeaderIndex(headers)
	for _, row := range rows {
		e := expenses.Expense{
			Expense:   utils.CSVVal(row[col["expense"]]),
			Cost:      *utils.CSVFloat(row[col["cost"]]),
			Frequency: utils.CSVVal(row[col["frequency"]]),
		}
		if i, ok := col["date"]; ok && i < len(row) {
			v := utils.CSVVal(row[i])
			if v != "" {
				e.Date = &v
			}
		}
		if i, ok := col["due_date"]; ok && i < len(row) {
			e.DueDate = utils.CSVInt(row[i])
		}
		repo.Create(userID, e)
	}
	logger.Log.Info().Str("file", file).Msg("seeded expenses")
}

func seedDebts(repo debts.Repository, userID int, file string) {
	headers, rows, ok := utils.ReadCSV(file)
	if !ok {
		return
	}
	col := utils.HeaderIndex(headers)
	for _, row := range rows {
		repo.Create(userID, debts.Debt{
			Bank:    utils.CSVVal(row[col["bank"]]),
			Type:    utils.CSVVal(row[col["type"]]),
			Name:    utils.CSVVal(row[col["name"]]),
			APY:     *utils.CSVFloat(row[col["apy"]]),
			Balance: *utils.CSVFloat(row[col["balance"]]),
		})
	}
	logger.Log.Info().Str("file", file).Msg("seeded debts")
}

func seedIncome(repo income.Repository, userID int, file string) {
	headers, rows, ok := utils.ReadCSV(file)
	if !ok {
		return
	}
	col := utils.HeaderIndex(headers)
	for _, row := range rows {
		inc := income.Income{
			Job:     utils.CSVVal(row[col["job"]]),
			PayType: utils.CSVVal(row[col["pay_type"]]),
		}
		if i, ok := col["pay_per_hour"]; ok && i < len(row) {
			inc.PayPerHour = utils.CSVFloat(row[i])
		}
		if i, ok := col["hour_per_day"]; ok && i < len(row) {
			inc.HourPerDay = utils.CSVFloat(row[i])
		}
		if i, ok := col["annual_salary"]; ok && i < len(row) {
			inc.AnnualSalary = utils.CSVFloat(row[i])
		}
		if i, ok := col["pay_frequency"]; ok && i < len(row) {
			v := utils.CSVVal(row[i])
			if v != "" {
				inc.PayFrequency = &v
			}
		}
		if i, ok := col["pay_day"]; ok && i < len(row) {
			inc.PayDay = utils.CSVInt(row[i])
		}
		repo.Create(userID, inc)
	}
	logger.Log.Info().Str("file", file).Msg("seeded income")
}

func seedLiquid(repo liquid.Repository, userID int, file string) {
	headers, rows, ok := utils.ReadCSV(file)
	if !ok {
		return
	}
	col := utils.HeaderIndex(headers)
	for _, row := range rows {
		repo.Create(userID, liquid.Liquid{
			Bank:    utils.CSVVal(row[col["bank"]]),
			Balance: *utils.CSVFloat(row[col["balance"]]),
		})
	}
	logger.Log.Info().Str("file", file).Msg("seeded liquid")
}

func seedInsights(repo insights.Repository, userID int) {
	response := `{
		"overview": "Your finances show a steady income from your Vet Assistant position at $19/hr, but high-interest credit card debt totaling over $17,000 is eating into your progress. With $78,000+ in total debt and only $1,864 in liquid savings, the priority should be aggressive debt payoff starting with your Capital One Savor card at 28.99% APR.",
		"health_score": 42,
		"strengths": [
			"Consistent biweekly income provides predictable cash flow for budgeting",
			"Student loan interest rates (8-10%) are lower than credit card rates, allowing strategic payoff ordering",
			"You have a zero-interest Chase Slate balance of $119.70 that can be cleared quickly for a psychological win"
		],
		"warnings": [
			"Capital One Savor card at 28.99% APR on $14,210 is generating roughly $343/month in interest alone",
			"Liquid savings of $1,864 covers less than 1 month of expenses — well below the 3-6 month emergency fund target",
			"Over 60% of monthly income is going toward debt minimums and fixed expenses, leaving minimal room for savings"
		],
		"advice": [
			{
				"title": "Eliminate the Capital One Savor First",
				"detail": "At 28.99% APR, your $14,210 Savor balance costs ~$343/month in interest. Every extra dollar toward this card has the highest return. Even an extra $100/month saves $1,200+ in interest over the payoff period."
			},
			{
				"title": "Clear the Chase Slate Immediately",
				"detail": "Your $119.70 Chase Slate balance has 0% APR — pay it off this month to eliminate one account entirely and simplify your debt picture."
			},
			{
				"title": "Build a $1,000 Starter Emergency Fund",
				"detail": "Before going all-in on debt, set aside $1,000 as a buffer so unexpected expenses don't force you back onto credit cards. You are $136 away from this milestone."
			}
		],
		"resources": [
			{
				"title": "Debt Avalanche Calculator",
				"description": "Model how targeting your 28.99% Capital One balance first minimizes total interest paid across all 9 accounts."
			},
			{
				"title": "r/personalfinance Prime Directive",
				"description": "Step-by-step flowchart for balancing emergency savings, high-interest debt payoff, and retirement contributions on a limited income."
			}
		]
	}`
	if err := repo.Save(userID, "insights", response); err != nil {
		logger.Log.Error().Err(err).Msg("seed insights error")
		return
	}
	logger.Log.Info().Msg("seeded insights cache for demo user")
}
