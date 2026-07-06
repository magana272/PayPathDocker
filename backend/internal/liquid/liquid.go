package liquid

type Liquid struct {
	ID      int     `json:"id" bson:"id"`
	UserID  int     `json:"-" bson:"user_id"`
	Bank    string  `json:"bank" bson:"bank"`
	Balance float64 `json:"balance" bson:"balance"`
}
