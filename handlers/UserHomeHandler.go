package handlers

import (
	"forum/config"
	"forum/middleware"
	"forum/models"
	"net/http"
	"strconv"
)

func HomeUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	userID, ok := r.Context().Value(middleware.UserIdKey).(int)
	if !ok {
		HandleError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Map des catégories valides
	validCats := map[string]bool{}
	for _, cat := range categories {
		validCats[strconv.Itoa(cat.IdCat)] = true
	}

	query := r.URL.Query()
	selectedCats := query["category"]
	activities := query["activity"]

	// Valider activities
	for _, activity := range activities {
		if activity != "liked" && activity != "owned" {
			HandleError(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if activity == "liked" {
			data.FilterLikes = true
		}
		if activity == "owned" {
			data.FilterMyPosts = true
		}
	}

	hasAll := false
	hasSpecificCat := false
	filteredCats := []string{}

	for _, cat := range selectedCats {
		if cat == "all" {
			hasAll = true
		} else {
			// Valider que la catégorie existe en base
			if !validCats[cat] {
				HandleError(w, "Bad Request", http.StatusBadRequest)
				return
			}
			hasSpecificCat = true
			filteredCats = append(filteredCats, cat)
		}
	}

	data.SelectedCategories = make(map[string]bool)
	if hasAll {
		data.SelectedCategories["all"] = true
	}
	for _, cat := range filteredCats {
		data.SelectedCategories[cat] = true
	}
	if !hasSpecificCat && !hasAll {
		data.SelectedCategories["all"] = true
	}

	hasActivity := data.FilterLikes || data.FilterMyPosts

	var posts []models.Post

	if hasAll || (!hasSpecificCat && !hasActivity) {
		posts, err = models.GetAllPosts()
		if err != nil {
			HandleError(w, "Error loading posts", http.StatusInternalServerError)
			return
		}
		data.SelectedCategories = map[string]bool{"all": true}
	} else {
		posts, err = models.GetFilteredPosts(userID, filteredCats, data.FilterLikes, data.FilterMyPosts)
		if err != nil {
			HandleError(w, "Error loading filtered posts", http.StatusInternalServerError)
			return
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
