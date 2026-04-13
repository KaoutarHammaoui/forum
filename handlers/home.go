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
	"strings"
	"time"
)

type Data struct {
	Posts      []models.Post
	Categories []models.Category
	Error      string
}

func HomeUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("TEST HOME USER")
	if r.Method == http.MethodPost {
		handleCreatePost(w, r)
		return
	}
	handleGetHome(w, r)
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categories := r.Form["categories"]

	if title == "" || content == "" || categories == nil {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if len(title) > 200 {
		http.Error(w, "Title too long", http.StatusBadRequest)
		return
	}

	imagePath, err := handleImageUpload(r)
	if err != nil {
		http.Error(w, "Image upload failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIdKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	post := models.Post{
		Title:   title,
		Content: content,
		UserId:  userID,
		Image:   imagePath,
	}

	postID, err := models.InsertPost(post)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// FIX: categories now sends Name (from HTML fix), so GetCategoryByName works correctly
	for _, catName := range categories {
		cat, err := models.GetCategoryByName(catName)
		if err != nil {
			continue
		}
		if err := models.InsertPostCategory(postID, cat.IdCat); err != nil {
			continue
		}
	}

	http.Redirect(w, r, "/homeUser", http.StatusSeeOther)
}

func handleGetHome(w http.ResponseWriter, r *http.Request) {
	data := Data{}

	categories, err := models.GetAllCategory()
	if err != nil {
		http.Error(w, "Error loading categories", http.StatusInternalServerError)
		return
	}
	data.Categories = categories

	posts, err := models.GetAllPosts()
	if err != nil {
		http.Error(w, "Error loading posts", 500)
		return
	}

	for i := range posts {
		posts[i].Comments, _ = models.GetCommentsByPost(posts[i].IdPost)
	}

	data.Posts = posts
	config.RenderTemplate(w, "homeUser.html", data)
}

func handleImageUpload(r *http.Request) (string, error) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		return "", nil
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowed[ext] {
		return "", fmt.Errorf("format non supporté (%s)", ext)
	}

	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(handler.Filename))
	dst := filepath.Join(uploadDir, uniqueName)

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}

	return dst, nil
}

func Home(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
    config.RenderTemplate(w, "home.html", nil)
}
