package handlers

import (
	"forum/config"
	"forum/middleware"
	"forum/models"
	"net/http"
)

func HomeUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := Data{}
	data.IsLogged = true
	data.SelectedCategories = map[string]bool{"all": true}

	categories, err := models.GetAllCategory()
	if err != nil {
		HandleError(w, "Error loading categories", http.StatusInternalServerError)
		return
	}
	data.Categories = categories

	var posts []models.Post
	userID, ok := r.Context().Value(middleware.UserIdKey).(int)
	if !ok {
		HandleError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		posts, err = models.GetAllPosts()
		if err != nil {
			HandleError(w, "Error loading posts", http.StatusInternalServerError)
			return
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			HandleError(w, "Bad Request", http.StatusBadRequest)
			return
		}

		selectedCats := r.Form["category"]
		data.SelectedCategories = make(map[string]bool, len(selectedCats))
		for _, category := range selectedCats {
			data.SelectedCategories[category] = true
		}
		for _, activity := range r.Form["activity"] {
			if activity == "liked" {
				data.FilterLikes = true
			}
			if activity == "owned" {
				data.FilterMyPosts = true
			}
		}

		hasSpecificCategory := false
		for _, category := range selectedCats {
			if category != "all" {
				hasSpecificCategory = true
				break
			}
		}

		if !hasSpecificCategory {
			data.SelectedCategories["all"] = true
		}

		if !data.FilterLikes && !data.FilterMyPosts && !hasSpecificCategory {
			posts, err = models.GetAllPosts()
			if err != nil {
				HandleError(w, "Error loading posts", http.StatusInternalServerError)
				return
			}
		} else {
			posts, err = models.GetFilteredPosts(userID, selectedCats, data.FilterLikes, data.FilterMyPosts)
			if err != nil {
				HandleError(w, "Error loading filtered posts", http.StatusInternalServerError)
				return
			}
		}
	}

	posts, err = GetInfoPosts(w, posts)
	if err != nil {
		HandleError(w, "Error loading posts", http.StatusInternalServerError)
		return
	}

	data.Posts = posts
	data.Action = "/homeUser"
	config.RenderTemplate(w, "homeUser.html", data)
}

func GetInfoPosts(w http.ResponseWriter, posts []models.Post) ([]models.Post, error) {
	for i, p := range posts {
		countlikes, err := models.CountLikeDislikeByPost(p.IdPost, "like")
		if err != nil {
			return nil, err
		}
		countdislikes, err := models.CountLikeDislikeByPost(p.IdPost, "dislike")
		if err != nil {
			return nil, err
		}
		comments, err := models.GetCommentsByPost(p.IdPost)
		if err != nil {
			return nil, err
		}
		posts[i].Likes = countlikes
		posts[i].Dislikes = countdislikes
		posts[i].Comments = comments
	}
	return posts, nil
}
