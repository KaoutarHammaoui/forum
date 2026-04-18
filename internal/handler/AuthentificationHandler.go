package handler

import (
	"forum/internal/config"
	models "forum/internal/model"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginData struct {
	Email         string
	EmailError    string
	PasswordError string
	HasErrors     bool
}

type RegisterData struct {
	Username      string
	Email         string
	UsernameError string
	EmailError    string
	PasswordError string
	HasErrors     bool
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		LoginHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		HandleError(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	data := LoginData{
		Email: r.URL.Query().Get("email"),
	}

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

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		RegisterHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		HandleError(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	config.RenderTemplate(w, "register.html", RegisterData{})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		HandleError(w, "bad request", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	user, err := models.GetUserByEmail(email)
	if err != nil {
		http.Redirect(w, r, "/login?error=email&email="+email, http.StatusSeeOther)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Redirect(w, r, "/login?error=password&email="+email, http.StatusSeeOther)
		return
	}

	if err := config.DeleteSessionsByUserID(user.ID); err != nil {
		HandleError(w, "internal error", http.StatusInternalServerError)
		return
	}

	token, err := config.InsertSession(user.ID)
	if err != nil {
		HandleError(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		HandleError(w, "bad request", http.StatusBadRequest)
		return
	}

	data := RegisterData{
		Username: strings.TrimSpace(r.FormValue("username")),
		Email:    strings.TrimSpace(r.FormValue("email")),
	}
	password := r.FormValue("password")

	if len(data.Username) < 3 {
		data.UsernameError = "Username must be at least 3 characters"
		data.HasErrors = true
	}

	if _, err := mail.ParseAddress(data.Email); err != nil {
		data.EmailError = "Enter a valid email address"
		data.HasErrors = true
	}

	if len(password) < 6 {
		data.PasswordError = "Password must be at least 6 characters"
		data.HasErrors = true
	}

	if exists, err := models.ExistsInColumn("username", data.Username); err == nil && exists {
		data.UsernameError = "Username is already taken"
		data.HasErrors = true
	}

	if exists, err := models.ExistsInColumn("email", data.Email); err == nil && exists {
		data.EmailError = "An account already exists with this email"
		data.HasErrors = true
	}

	if data.HasErrors {
		config.RenderTemplate(w, "register.html", data)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		HandleError(w, "internal error", http.StatusInternalServerError)
		return
	}

	userID, err := models.InsertUser(models.User{
		Username: data.Username,
		Email:    data.Email,
		Password: string(hash),
	})
	if err != nil {
		HandleError(w, "internal error", http.StatusInternalServerError)
		return
	}

	token, err := config.InsertSession(int(userID))
	if err != nil {
		HandleError(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogOUT(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("token")
	if err == nil {
		_ = config.DeleteSessionByToken(cookie.Value)
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
