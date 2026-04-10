package handlers

import (
	"forum/config"
	"forum/middleware"
	"forum/models"
	"net/http"
	"strconv"
	"strings"
)

func HomeHAndler(w http.ResponseWriter, r *http.Request) {
	categories, err := models.GetAllCategories()
	if err != nil {
		http.Error(w, "categories not found", http.StatusInternalServerError)
		return
	}

	data := HomeData{
		Categories: categories,
	}

	if r.Method == http.MethodPost {
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		selectedCategories := r.Form["category"]

		data.Title = title
		data.Content = content

		if title == "" {
			data.TitleError = "title is required"
		}
		if content == "" {
			data.ContentError = "content is required"
		}
		if len(selectedCategories) == 0 {
			data.CategoryError = "choose at least one category"
		}

		if data.TitleError == "" && data.ContentError == "" && data.CategoryError == "" {
			userID, ok := r.Context().Value(middleware.UserIdKey).(int)
			if !ok {
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}

			postID, err := models.InsertPost(models.Post{
				Title:   title,
				Content: content,
				UserId:  userID,
			})
			if err != nil {
				http.Error(w, "post not inserted", http.StatusInternalServerError)
				return
			}

			for _, categoryID := range selectedCategories {
				id, err := strconv.Atoi(categoryID)
				if err != nil {
					http.Error(w, "wrong category", http.StatusBadRequest)
					return
				}

				err = models.InsertPostCategory(int(postID), id)
				if err != nil {
					http.Error(w, "post category not inserted", http.StatusInternalServerError)
					return
				}
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	} else if r.Method != http.MethodGet {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tmpl := config.GetTemplate("home.html")
	if tmpl == nil {
		http.Error(w, "not found", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "home.html", data)
}
