package handlers

import (
	"forum/config"
	"forum/models"
	"net/http"
	"strconv"
)

func HomeUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	handleGetHomeUser(w, r)
	
}

func handleGetHomeUser(w http.ResponseWriter, r *http.Request) {
	data := Data{}
	data.IsLogged=true
	categories, err := models.GetAllCategory()

	if err != nil {
		HandleError(w, "Error loading categories", http.StatusInternalServerError)
		return
	}
	data.Categories = categories

	if r.Method == http.MethodGet {
		posts, err := models.GetAllPosts()
		if err != nil {
			HandleError(w, "Error loading posts", http.StatusInternalServerError)
			return
		}
		for i, p := range posts {
			comments, err := models.FetchComment(p.IdPost)
			if err != nil {
				continue
			}
			posts[i].Comments = comments
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
				HandleError(w, "Error loading posts", http.StatusInternalServerError)
				return
			}

			//fetching post comments
			for i, p := range posts {
				comments, err := models.FetchComment(p.IdPost)
				if err != nil {
					continue
				}
				posts[i].Comments = comments
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
