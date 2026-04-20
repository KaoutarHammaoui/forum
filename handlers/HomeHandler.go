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
	if r.Method != http.MethodGet {
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
	data.SelectedCategories = map[string]bool{"all": true}
	data.Action = "/"

	// Map des catégories valides
	validCats := map[string]bool{}
	for _, cat := range categories {
		validCats[strconv.Itoa(cat.IdCat)] = true
	}

	selectedCats := r.URL.Query()["category"]

	hasAll := false
	hasSpecificCat := false
	filteredCats := []string{}

	for _, cat := range selectedCats {
		if cat == "all" {
			hasAll = true
		} else {
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

	var posts []models.Post

	if hasAll || !hasSpecificCat {
		data.SelectedCategories = map[string]bool{"all": true}
		posts, err = models.GetAllPosts()
		if err != nil {
			HandleError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		postMap := map[int]models.Post{}
		for _, catStr := range filteredCats {
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

	posts, err = GetInfoPosts(w, posts)
	if err != nil {
		HandleError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data.Posts = posts
	config.RenderTemplate(w, "home.html", data)
}