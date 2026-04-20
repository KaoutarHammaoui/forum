package handlers

import (
	"fmt"
	"forum/config"
	"forum/middleware"
	"forum/models"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	handleCreatePost(w, r)
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		renderCreateError(w, r, "Invalid form data.")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categories := r.Form["categories"]

	// Title
	if title == "" {
		renderCreateError(w, r, "Title is required.")
		return
	}
	if len(title) > 200 {
		renderCreateError(w, r, "Title must not exceed 200 characters.")
		return
	}

	// Content
	if content == "" {
		renderCreateError(w, r, "Content is required.")
		return
	}
	if len(content) > 1000 {
		renderCreateError(w, r, "Content must not exceed 1000 characters.")
		return
	}

	// Categories
	if len(categories) == 0 {
		renderCreateError(w, r, "Please select at least one category.")
		return
	}

	// Image
	file, handler, err := r.FormFile("image")
	if err != nil {
		renderCreateError(w, r, "An image is required.")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		renderCreateError(w, r, "Image must be a JPG or JPEG or PNG file.")
		return
	}

	if handler.Size > 5<<20 {
		renderCreateError(w, r, "Image must not exceed 5 MB.")
		return
	}

	// Save image
	imagePath, err := saveUploadedFile(file, handler.Filename, ext)
	if err != nil {
		renderCreateError(w, r, "Failed to save image. Please try again.")
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
		Image:   imagePath,
	}
	postID, err := models.InsertPost(post)
	if err != nil {
		renderCreateError(w, r, "Failed to create post. Please try again.")
		return
	}

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

func saveUploadedFile(file multipart.File, originalName, ext string) (string, error) {
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}
	uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(originalName))
	dst := filepath.Join(uploadDir, uniqueName)

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}
	return strings.ReplaceAll(dst, "\\", "/"), nil
}
