package middleware

import (
	"context"
	"forum/internal/service"
	"net/http"
	"time"
)

func AuthMiddleWare(next http.HandlerFunc, service *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		session, err := service.Authorization.GetSessionByToken(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		user, err := service.Authorization.GetUserByID(session.UserID)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next(w, r.WithContext(ctx))
	}
}
