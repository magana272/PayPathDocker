package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "user_id"

func ContextWithUserID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func UserID(r *http.Request) int {
	return r.Context().Value(userIDKey).(int)
}

type Authenticator interface {
	Authenticate(token string) (int, error)
}

func RequireAuth(a Authenticator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing token", 401)
			return
		}
		uid, err := a.Authenticate(strings.TrimPrefix(authHeader, "Bearer "))
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}
		next.ServeHTTP(w, r.WithContext(ContextWithUserID(r.Context(), uid)))
	})
}
