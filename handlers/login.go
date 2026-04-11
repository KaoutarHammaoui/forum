package handlers

import (
	"forum/config"
	"forum/models"
	"net/http"
	"strings"
	"time"
	"golang.org/x/crypto/bcrypt"
)

func LoginH(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}
	data := Login{}
	switch r.URL.Query().Get("error") {
	case "email":
		data.EmailError = "No account found with this email"
		data.HasErrors = true
	case "password":
		data.PasswordError = "incorrect password"
		data.HasErrors = true
	}
	config.RenderTemplate(w, "login.html", data)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return

	}
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	data := Login{Email: email}

	user, err := models.GetUserByEmail(email)
	if err != nil {
		data.EmailError = "no account with this email"
		data.HasErrors = true
		http.Redirect(w, r, "/login?error=email", http.StatusSeeOther)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		data.PasswordError = "incorrect password"
		data.HasErrors = true
		http.Redirect(w, r, "/login?error=password", http.StatusSeeOther)
		return
	}

	token, err := models.InsertSession(user.ID)
	if err != nil {
		http.Error(w, "wrong", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Path:     "/homeUser",
	})
	http.Redirect(w, r, "/homeUser", http.StatusSeeOther)
}

func LogOUT(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")

	if err == nil {
		models.DeleteSessionByToken(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
