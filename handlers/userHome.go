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

	if r.Method == http.MethodPost {
		handleCreatePost(w, r)
		return
	}
	handleGetHome(w, r)
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {

	// 1. Parse du formulaire multipart (10 MB max)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 2. Récupération et validation des champs obligatoires
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categories := r.Form["categories"]

	if title == "" || content == "" || categories == nil {
		//
		return
	}
	if len(title) > 200 {
		//
		return
	}

	// 3. Upload image (optionnel)
	imagePath, err := handleImageUpload(r)
	if err != nil {
		//
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
		//
	}
	// 5. Liaison post ↔ catégories (récupération de l'id par le nom)
	for _, catName := range categories {
		cat, err := models.GetCategoryByName(catName)
		if err != nil {
			continue
		}
		if err := models.InsertPostCategory(postID, cat.IdCat); err != nil {
			continue
		}
	}

	// 6. Redirect POST → GET (pattern PRG)
	http.Redirect(w, r, "/homeUser", http.StatusSeeOther)
}

// ─────────────────────────────────────────────────────────────────────────────

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
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}
	data.Posts = posts

	config.RenderTemplate(w, "homeUser.html", data)
}

// ─────────────────────────────────────────────────────────────────────────────

// handleImageUpload gère l'upload et retourne le chemin relatif, ou "" si absent.
func handleImageUpload(r *http.Request) (string, error) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		// champ absent ou vide → pas d'image, pas d'erreur
		return "", nil
	}
	defer file.Close()

	// Extensions autorisées
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowed[ext] {
		return "", fmt.Errorf("format non supporté (%s)", ext)
	}

	// Dossier uploads
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	// Nom unique pour éviter les collisions
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
