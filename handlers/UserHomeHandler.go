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

	data := Data{}
	data.IsLogged = true

	categories, err := models.GetAllCategory()
	if err != nil {
		HandleError(w, "Error loading categories", http.StatusInternalServerError)
		return
	}
	data.Categories = categories

	var posts []models.Post

	if r.Method == http.MethodGet {
		posts, err = models.GetAllPosts()
		if err != nil {
			HandleError(w, "Error loading posts", http.StatusInternalServerError)
			return
		}
	} else {
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
			posts, err = models.GetAllPosts()
			if err != nil {
				HandleError(w, "Error loading posts", http.StatusInternalServerError)
				return
			}
		} else {
			postMap := map[int]models.Post{}
			for _, catStr := range selectedCats {
				catId, err := strconv.Atoi(catStr)
				if err != nil {
					continue
				}
				catPosts, err := models.GetPostsByCategory(catId)
				if err != nil {
					continue
				}
				for _, p := range catPosts {
					postMap[p.IdPost] = p
				}
			}
			for _, p := range postMap {
				posts = append(posts, p)
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
