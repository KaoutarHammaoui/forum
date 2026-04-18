package handler

import (
	"forum/internal/config"
	models "forum/internal/model"

	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginData struct {
	EmailError    string
	PasswordError string
	HasErrors     bool
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	data := LoginData{}

	switch r.URL.Query().Get("error") {
	case "email":
		data.EmailError = "No account found with this email"
		data.HasErrors = true
	case "password":
		data.PasswordError = "Incorrect password"
		data.HasErrors = true
	}

	config.RenderTemplate(w, "login.html", data)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form (safe practice)
	if err := r.ParseForm(); err != nil {
		HandleError(w, "bad request", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	user, err := models.GetUserByEmail(email)
	if err != nil {
		http.Redirect(w, r, "/login?error=email", http.StatusSeeOther)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Redirect(w, r, "/login?error=password", http.StatusSeeOther)
		return
	}

	// Delete old sessions (optional but good practice)
	config.DeleteSessionsByUserID(user.ID)

	token, err := config.InsertSession(user.ID)
	if err != nil {
		HandleError(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		// Secure: true, // enable in HTTPS
		// SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/homeUser", http.StatusSeeOther)
}

func LogOUT(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err == nil {
		config.DeleteSessionByToken(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}