package middleware

import (
	"context"
	"net/http"
	"time"

	"forum/models"
)

type contextKey string

const UserIdKey contextKey = "userID"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session, err := models.GetSessionByToken(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   "",
				Expires: time.Now().Add(-time.Hour),
				HttpOnly: true,
				Path:    "/",
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		cts := context.WithValue(r.Context(), UserIdKey, session.UserId)
		next(w, r.WithContext(cts))
	}
}
