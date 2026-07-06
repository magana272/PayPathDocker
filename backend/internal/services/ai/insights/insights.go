package insights

import "time"

type InsightsCache struct {
	UserID    int       `bson:"user_id"`
	Topic     string    `bson:"topic"`
	Response  string    `bson:"response"`
	CreatedAt time.Time `bson:"created_at"`
}

type Advice struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type Resource struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Insights struct {
	Overview    string     `json:"overview"`
	HealthScore int        `json:"health_score"`
	Strengths   []string   `json:"strengths"`
	Warnings    []string   `json:"warnings"`
	Advice      []Advice   `json:"advice"`
	Resources   []Resource `json:"resources"`
}
