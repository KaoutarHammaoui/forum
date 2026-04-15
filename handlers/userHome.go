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
    if r.Method != http.MethodGet && r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // POST : soit création de post, soit filtrage par catégorie
    if r.Method == http.MethodPost {
        // Si le form contient "title", c'est une création de post
        if r.FormValue("title") != "" {
            handleCreatePost(w, r)
            return
        }
        // Sinon c'est un filtre par catégorie
        handleGetHomeUser(w, r)
        return
    }

    handleGetHomeUser(w, r)
}

func handleGetHomeUser(w http.ResponseWriter, r *http.Request) {
    data := Data{}

    categories, err := models.GetAllCategory()
    if err != nil {
        http.Error(w, "Error loading categories", http.StatusInternalServerError)
        return
    }
    data.Categories = categories

    if r.Method == http.MethodGet {
        posts, err := models.GetAllPosts()
        if err != nil {
            http.Error(w, "Error loading posts", http.StatusInternalServerError)
            return
        }
        data.Posts = posts

    } else if r.Method == http.MethodPost {
        r.ParseForm()
        selectedCats := r.Form["category"]

        isAll := false
        for _, v := range selectedCats {
            if v == "all" {
                isAll = true
                break
            }
        }

        if isAll || len(selectedCats) == 0 {
            posts, err := models.GetAllPosts()
            if err != nil {
                http.Error(w, "Error loading posts", http.StatusInternalServerError)
                return
            }
            data.Posts = posts
        } else {
            postMap := map[int]models.Post{}
            for _, catStr := range selectedCats {
                catId, err := strconv.Atoi(catStr)
                if err != nil {
                    continue
                }
                posts, err := models.GetPostsByCategory(catId)
                if err != nil {
                    continue
                }
                for _, p := range posts {
                    postMap[p.IdPost] = p
                }
            }
            for _, p := range postMap {
                data.Posts = append(data.Posts, p)
            }
        }
    }

    data.Action = "/homeUser"
    config.RenderTemplate(w, "homeUser.html", data)
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
		//
		return
	}
	if len(title) > 200 {
		//
		return
	}
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
	for _, catIDStr := range categories {
		catID, err := strconv.Atoi(catIDStr)

		if err != nil {
			return
		}
		if err := models.InsertPostCategory(postID, catID); err != nil {
			return
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
	fmt.Println(dst)

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}
	return strings.ReplaceAll(filepath.Join(uploadDir, uniqueName), "\\", "/"), nil
}
