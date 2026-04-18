package handlers

import (
	"forum/config"
	"forum/middleware"
	"forum/models"
	"net/http"
	"strconv"
)

func Home(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == http.MethodGet {

		posts, err := models.GetAllPosts()
		if err != nil {
			HandleError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		for i, p := range posts {
			countlikes, err := models.CountLikeDislikeByPost(p.IdPost, "like")
			if err != nil {

				return
			}
			posts[i].Likes = countlikes

			Countdislikes, err := models.CountLikeDislikeByPost(p.IdPost, "dislike")
			if err != nil {
				HandleError(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			posts[i].Dislikes = Countdislikes
			comments, err := models.FetchComment(p.IdPost)
			if err != nil {
				HandleError(w, "Internal Server Error", http.StatusInternalServerError)
				return
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
				HandleError(w, "Internal Server Error", http.StatusInternalServerError)
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
	} else {
		HandleError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	data.Action = "/"
	config.RenderTemplate(w, "home.html", data)
}
