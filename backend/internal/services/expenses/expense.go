package expenses

type DateException struct {
	OriginalDate string `json:"original_date" bson:"original_date"`
	NewDate      string `json:"new_date" bson:"new_date"`
}

type Expense struct {
	ID         int             `json:"id" bson:"id"`
	UserID     int             `json:"-" bson:"user_id"`
	Expense    string          `json:"expense" bson:"expense"`
	Cost       float64         `json:"cost" bson:"cost"`
	Date       *string         `json:"date" bson:"date,omitempty"`
	DueDate    *int            `json:"due_date" bson:"due_date,omitempty"`
	Frequency  string          `json:"frequency" bson:"frequency"`
	Exceptions []DateException `json:"exceptions" bson:"exceptions,omitempty"`
}
