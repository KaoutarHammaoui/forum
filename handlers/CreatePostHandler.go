package handlers

import (
	"forum/config"
	"forum/middleware"
	"forum/models"
	"net/http"
	"strconv"
	"strings"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	handleCreatePost(w, r)
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderCreateError(w, r, err.Error())
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categories := r.Form["categories"]

	// Title validation
	if title == "" {
		renderCreateError(w, r, "Title is required.")
		return
	}
	if len(title) > 200 {
		renderCreateError(w, r, "Title must not exceed 200 characters.")
		return
	}

	// Content validation
	if content == "" {
		renderCreateError(w, r, "Content is required.")
		return
	}
	if len(content) > 1000 {
		renderCreateError(w, r, "Content must not exceed 1000 characters.")
		return
	}

	// Categories validation
	if len(categories) == 0 {
		renderCreateError(w, r, "Please select at least one category.")
		return
	}
	// User
	userID, ok := r.Context().Value(middleware.UserIdKey).(int)
	if !ok {
		HandleError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Insert post
	post := models.Post{
		Title:   title,
		Content: content,
		UserId:  userID,
	}

	postID, err := models.InsertPost(post)
	if err != nil {
		renderCreateError(w, r, "Failed to create post. Please try again.")
		return
	}

	// Attach categories
	for _, catIDStr := range categories {
		catID, err := strconv.Atoi(catIDStr)
		if err != nil {
			renderCreateError(w, r, "Invalid category selected.")
			return
		}
		if err := models.InsertPostCategory(postID, catID); err != nil {
			renderCreateError(w, r, "Failed to attach category. Please try again.")
			return
		}
	}

	http.Redirect(w, r, "/homeUser", http.StatusSeeOther)
}

func renderCreateError(w http.ResponseWriter, r *http.Request, msg string) {
	data := Data{Error: msg}

	categories, err := models.GetAllCategory()
	if err == nil {
		data.Categories = categories
	}

	posts, err := models.GetAllPosts()
	if err == nil {
		data.Posts = posts
	}

	data.Action = "/homeUser"
	w.WriteHeader(404)
	config.RenderTemplate(w, "homeUser.html", data)
}
