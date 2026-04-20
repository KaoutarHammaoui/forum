package handlers

import "forum/models"

type Login struct {
	Email         string
	EmailError    string
	PasswordError string
	HasErrors     bool
}

type Data struct {
	Posts        []models.Post
	Categories   []models.Category
	Error        string
	Action       string
	IsLogged     bool
}

type RegistrationData struct {
	Username      string
	Email         string
	Password      string
	UsernameError string
	EmailError    string
	PasswordError string
	HasErrors     bool
}
