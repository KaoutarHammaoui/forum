package handlers

import (
	"fmt"
	"forum/config"
	"forum/middleware"
	"forum/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func HomeUser(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		handleGetHomeUser(w, r, "")

	case http.MethodPost:
		handleCreatePost(w, r)

	default:
		Error(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleGetHomeUser(w http.ResponseWriter, r *http.Request, errMsg string) {

	data := Data{
		ErrorMsg: errMsg,
	}

	categories, err := models.GetAllCategory()
	if err != nil {
		Error(w, http.StatusInternalServerError, "Error loading categories")
		return
	}
	data.Categories = categories

	posts, err := models.GetAllPosts()
	if err != nil {
		Error(w, http.StatusInternalServerError, "Error loading posts")
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		selectedCats := r.Form["categories"]

		if len(selectedCats) > 0 && !contains(selectedCats, "all") {
			postMap := make(map[int]models.Post)

			for _, catStr := range selectedCats {
				catID, err := strconv.Atoi(catStr)
				if err != nil {
					continue
				}

				catPosts, err := models.GetPostsByCategory(catID)
				if err != nil {
					continue
				}

				for _, p := range catPosts {
					postMap[p.IdPost] = p
				}
			}

			posts = []models.Post{}
			for _, p := range postMap {
				posts = append(posts, p)
			}
		}
	}

	data.Posts = posts
	data.Action = "/homeUser"

	config.RenderTemplate(w, "homeUser.html", data)
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categories := r.Form["categories"]

	// ✅ VALIDATION SIMPLE
	if title == "" || content == "" || len(categories) == 0 {
		handleGetHomeUser(w, r, "All fields are required")
		return
	}

	if len(title) > 200 {
		handleGetHomeUser(w, r, "Title must be less than 200 characters")
		return
	}

	imagePath, imgErr := handleImageUpload(r)
	if imgErr != nil {
		handleGetHomeUser(w, r, imgErr.Error())
		return
	}

	userID, ok := r.Context().Value(middleware.UserIdKey).(int)
	if !ok {
		Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	postID, err := models.InsertPost(models.Post{
		Title:   title,
		Content: content,
		UserId:  userID,
		Image:   imagePath,
	})
	if err != nil {
		Error(w, http.StatusInternalServerError, "Error creating post")
		return
	}

	for _, catIDStr := range categories {
		if catID, err := strconv.Atoi(catIDStr); err == nil {
			models.InsertPostCategory(postID, catID)
		}
	}

	http.Redirect(w, r, "/homeUser", http.StatusSeeOther)
}

func handleImageUpload(r *http.Request) (string, error) {

	file, handler, err := r.FormFile("image")
	if err != nil {
		return "", nil
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	allowed := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true,
	}

	if !allowed[ext] {
		return "", fmt.Errorf("unsupported image format")
	}

	uploadDir := "uploads"
	os.MkdirAll(uploadDir, 0755)

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(handler.Filename))
	path := filepath.Join(uploadDir, filename)

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(path, "\\", "/"), nil
}

func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
