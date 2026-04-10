package handlers

import "forum/models"

type Login struct {
	Email         string
	EmailError    string
	PasswordError string
	HasErrors     bool
}

type HomeData struct {
	Title         string
	Content       string
	Categories    []models.Category
	TitleError    string
	ContentError  string
	CategoryError string
}
