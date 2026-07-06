package auth

import "time"

type User struct {
	ID       int    `json:"id" bson:"id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"-" bson:"password"`
	Name     string `json:"name" bson:"name"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *User  `json:"user,omitempty"`
	Token string `json:"token"`
}

type RevokedToken struct {
	Token     string    `bson:"token"`
	RevokedAt time.Time `bson:"revoked_at"`
}
