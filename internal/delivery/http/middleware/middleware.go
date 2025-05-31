package router

import (
	"LostAndFound/internal/auth"
	"context"
	"net/http"
	"strings"
)

const (
	ctxUserIDKey string = "userID"
	ctxRoleKey   string = "role"
)

func AuthMiddleware(tokenManager *auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			fields := strings.Split(authHeader, " ")
			if len(fields) != 2 || fields[0] != "Bearer" {
				http.Error(w, "invalid token format", http.StatusUnauthorized)
				return
			}

			isBlacklisted, err := tokenManager.CacheRepo.IsTokenBlacklisted(r.Context(), fields[1])
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			if isBlacklisted {
				http.Error(w, "token expired or logged out", http.StatusUnauthorized)
				return
			}

			claims, err := tokenManager.Parse(fields[1])
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			userID := claims.UserID
			role := claims.Role
			ctx := context.WithValue(r.Context(), ctxUserIDKey, userID)
			ctx = context.WithValue(ctx, ctxRoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) string {
	id, _ := ctx.Value(ctxUserIDKey).(string)
	return id
}

func GetUserRole(ctx context.Context) string {
	role, _ := ctx.Value(ctxRoleKey).(string)
	return role
}

func AdminOnlyMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := GetUserRole(r.Context())
			if role != "admin" {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
