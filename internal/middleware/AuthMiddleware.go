package middleware

import (
	"context"
	"forum/internal/config"
	"net/http"
)

type contextKey string

const UserIdKey contextKey = "userID"

// IsAuthenticated is your helper to check status inside Handlers
func IsAuthenticated(r *http.Request) (int, bool) {
    userID, ok := r.Context().Value(UserIdKey).(int)
    if !ok || userID == 0 {
        return 0, false
    }
    return userID, true
}

// CheckUserContext: Allows access to everyone, but identifies users if they have a cookie
func CheckUserContext(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("token")
        if err != nil {
            // No cookie? No problem. Just move to the next handler.
            next(w, r)
            return
        }

        session, err := config.GetSessionByToken(cookie.Value)
        if err != nil {
            // Invalid/Expired session? Just move on (optionally clear the bad cookie)
            next(w, r)
            return
        }

        // User is valid! Add ID to context so we know who they are.
        ctx := context.WithValue(r.Context(), UserIdKey, session.UserId)
        next(w, r.WithContext(ctx))
    }
}

// AuthMiddleware: Strict check for protected actions (create post, etc.)
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        _, loggedIn := IsAuthenticated(r)
        if !loggedIn {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        next(w, r)
    }
}