package debts

type Debt struct {
	ID      int     `json:"id" bson:"id"`
	UserID  int     `json:"-" bson:"user_id"`
	Bank    string  `json:"bank" bson:"bank"`
	Type    string  `json:"type" bson:"type"`
	Name    string  `json:"name" bson:"name"`
	APY     float64 `json:"apy" bson:"apy"`
	Balance float64 `json:"balance" bson:"balance"`
}
