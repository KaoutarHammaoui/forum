package handlers

import (
	"forum/config"
	"forum/middleware"
	"forum/models"
	"net/http"
	"strconv"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		HandleError(w, "Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		HandleError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := middleware.GetSession(r)
	if err == nil && session != nil {
		http.Redirect(w, r, "/homeUser", http.StatusSeeOther)
		return
	}

	categories, err := models.GetAllCategory()
	if err != nil {
		HandleError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := Data{}
	data.Categories = categories

	// ================= GET =================
	if r.Method == http.MethodGet {
		posts, err := models.GetAllPosts()
		if err != nil {
			HandleError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		posts, err = enrichPosts(posts)
		if err != nil {
			HandleError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data.Posts = posts
	}

	// ================= POST (FILTER) =================
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			HandleError(w, "Bad Request", http.StatusBadRequest)
			return
		}

		selectedCats := r.Form["category"]

		isAll := false
		for _, c := range selectedCats {
			if c == "all" {
				isAll = true
				break
			}
		}

		var posts []models.Post

		if isAll || len(selectedCats) == 0 {
			posts, err = models.GetAllPosts()
			if err != nil {
				HandleError(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else {
			postMap := map[int]models.Post{}

			for _, catStr := range selectedCats {
				catId, err := strconv.Atoi(catStr)
				if err != nil {
					HandleError(w, "Invalid category id", http.StatusBadRequest)
					return
				}

				catPosts, err := models.GetPostsByCategory(catId)
				if err != nil {
					HandleError(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				for _, p := range catPosts {
					postMap[p.IdPost] = p
				}
			}

			for _, p := range postMap {
				posts = append(posts, p)
			}
		}

		posts, err = enrichPosts(posts)
		if err != nil {
			HandleError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data.Posts = posts
	}

	data.Action = "/"
	config.RenderTemplate(w, "home.html", data)
}

// ================= HELPER =================
func enrichPosts(posts []models.Post) ([]models.Post, error) {
	for i, p := range posts {

		likes, err := models.CountLikeDislikeByPost(p.IdPost, "like")
		if err != nil {
			return nil, err
		}
		posts[i].Likes = likes

		dislikes, err := models.CountLikeDislikeByPost(p.IdPost, "dislike")
		if err != nil {
			return nil, err
		}
		posts[i].Dislikes = dislikes

		comments, err := models.GetCommentsByPost(p.IdPost)
		if err != nil {
			return nil, err
		}
		posts[i].Comments = comments
	}

	return posts, nil
}