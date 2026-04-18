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
	data.IsLogged = true
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
			countlikes, err := models.CountLikeDislikeByPost(p.IdPost, "like")
			if err != nil {
				HandleError(w, "Error loading posts", http.StatusInternalServerError)
				return
			}
			countdislike, err := models.CountLikeDislikeByPost(p.IdPost, "dislike")
			if err != nil {
				HandleError(w, "Error loading posts", http.StatusInternalServerError)
				return
			}
			comments, err := models.GetCommentsByPost(p.IdPost)
			if err != nil {
				HandleError(w, "Error loading posts", http.StatusInternalServerError)
				return
			}
			posts[i].Likes = countlikes
			posts[i].Dislikes = countdislike
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
				comments, err := models.GetCommentsByPost(p.IdPost)
				if err != nil {
					HandleError(w, "Error loading posts", http.StatusInternalServerError)
					return
				}
				countlikes, err := models.CountLikeDislikeByPost(p.IdPost, "like")
				if err != nil {
					HandleError(w, "Error loading posts", http.StatusInternalServerError)
					return
				}
				countdislikes, err := models.CountLikeDislikeByPost(p.IdPost, "dislike")
				if err != nil {
					HandleError(w, "Error loading posts", http.StatusInternalServerError)
					return
				}
				posts[i].Likes = countlikes
				posts[i].Dislikes = countdislikes
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
			for i, p := range data.Posts {
				comments, _ := models.GetCommentsByPost(p.IdPost)
				countlikes, _ := models.CountLikeDislikeByPost(p.IdPost, "like")
				countdislikes, _ := models.CountLikeDislikeByPost(p.IdPost, "dislike")
				data.Posts[i].Comments = comments
				data.Posts[i].Likes = countlikes
				data.Posts[i].Dislikes = countdislikes
			}
		}
	}
	data.Action = "/homeUser"
	config.RenderTemplate(w, "homeUser.html", data)
}
