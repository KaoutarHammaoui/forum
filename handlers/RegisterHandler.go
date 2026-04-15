package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"forum/config"
	"forum/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Handle GET request to show the registration form
	if r.Method != http.MethodGet {
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the pre-parsed template from cache
	tmpl := config.GetTemplate("register.html")
	if tmpl == nil {
		HandleError(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.Execute(w, nil)
	if err != nil {
		HandleError(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func DoRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := strings.TrimSpace(r.PostFormValue("username"))
	email := strings.TrimSpace(r.PostFormValue("email"))
	password := r.PostFormValue("password")

	// 1. Validate Input using the new struct
	data := ValidateRegistrationInput(username, email, password)

	if data.HasErrors {
		tmpl := config.GetTemplate("register.html")
		if tmpl == nil {
			HandleError(w, "Template not found", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// Pass the struct 'data' directly to the template
		tmpl.Execute(w, data)
		return
	}

	// 2. Hash and Save (Standard logic)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	newUser := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if _, err := models.InsertUser(newUser); err != nil {
		HandleError(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ValidateRegistrationInput(username, emailAddr, password string) RegistrationData {
	data := RegistrationData{
		Username: username,
		Email:    emailAddr,
	}

	// 1. Username Regex
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9]{3,20}$`)
	if !nameRegex.MatchString(username) {
		data.UsernameError = "Username must be 3-20 characters long and alphanumeric"
		data.HasErrors = true
	} else {
		// Only check DB if format is valid
		if exists, _ := models.ExistsInColumn("username", username); exists {
			data.UsernameError = "Username already taken"
			data.HasErrors = true
		}
	}

	// 2. Email Validation
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(strings.ToLower(emailAddr)) {
		data.EmailError = "Invalid email format"
		data.HasErrors = true
	} else {
		if exists, _ := models.ExistsInColumn("email", emailAddr); exists {
			data.EmailError = "Email already registered"
			data.HasErrors = true
		}
	}

	// 3. Password Validation
	if len(password) < 6 {
		data.PasswordError = "Password must be at least 6 characters long"
		data.HasErrors = true
	}

	return data
}
