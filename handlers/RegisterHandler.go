package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"forum/config"
	"forum/middleware"
	"forum/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.GetSession(r)
	if err == nil && session != nil {
		http.Redirect(w, r, "/homeUser", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodGet {
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config.RenderTemplate(w, "register.html", nil)
}

func DoRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username := strings.TrimSpace(r.PostFormValue("username"))
	email := strings.ToLower(strings.TrimSpace(r.PostFormValue("email")))
	password := r.PostFormValue("password")

	data := ValidateRegistrationInput(username, email, password)

	if data.HasErrors {
		config.RenderTemplate(w, "register.html", data)
		return
	}

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errHash != nil {
		HandleError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	newUser := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	if _, err := models.InsertUser(newUser); err != nil {
		HandleError(w, "Database Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ValidateRegistrationInput(username, emailAddr, password string) RegistrationData {
	data := RegistrationData{
		Username: username,
		Email:    emailAddr,
	}
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9]{3,20}$`)
	if !nameRegex.MatchString(username) {
		data.UsernameError = "Username must be 3-20 characters long and alphanumeric"
		data.HasErrors = true
	} else {
		exists, err := models.ExistsInColumn("username", username)
		if err != nil {
			data.HasErrors = true
			data.EmailError = "Internal server error"
			return data
		}

		if exists {
			data.EmailError = "username already registered"
			data.HasErrors = true
		}

	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(emailAddr) {
		data.EmailError = "Invalid email format"
		data.HasErrors = true
	} else {
		exists, err := models.ExistsInColumn("email", emailAddr)
		if err != nil {
			data.HasErrors = true
			data.EmailError = "Internal server error"
			return data
		}

		if exists {
			data.EmailError = "Email already registered"
			data.HasErrors = true
		}
	}
	if len(password) < 6 {
		data.PasswordError = "Password must be at least 6 characters long"
		data.HasErrors = true
	}
	return data
}
