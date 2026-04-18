package handler

import (
	"forum/internal/config"
	"forum/internal/middleware"
	models "forum/internal/model"
	"net/http"
	"strconv"
	"strings"
)

type Data struct {
	Title              string
	IsLoggedIn         bool
	UserID             int
	Categories         []models.Category
	Posts              []models.Post
	SelectedCategoryID int
	SelectedView       string
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, loggedIn := middleware.IsAuthenticated(r)

	categories, err := models.GetAllCategory()
	if err != nil {
		HandleError(w, "could not load categories", http.StatusInternalServerError)
		return
	}

	posts, selectedCategoryID, selectedView, err := loadPosts(r, userID, loggedIn)
	if err != nil {
		HandleError(w, "could not load posts", http.StatusInternalServerError)
		return
	}

	data := Data{
		Title:              "Home",
		IsLoggedIn:         loggedIn,
		UserID:             userID,
		Categories:         categories,
		Posts:              posts,
		SelectedCategoryID: selectedCategoryID,
		SelectedView:       selectedView,
	}

	config.RenderTemplate(w, "index.html", data)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, loggedIn := middleware.IsAuthenticated(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		HandleError(w, "bad request", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categoryID, _ := strconv.Atoi(r.FormValue("category_id"))

	if title == "" || content == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	postID, err := models.InsertPost(models.Post{
		Title:   title,
		Content: content,
		UserId:  userID,
	})
	if err != nil {
		HandleError(w, "could not create post", http.StatusInternalServerError)
		return
	}

	if categoryID > 0 {
		if err := models.InsertPostCategory(postID, categoryID); err != nil {
			HandleError(w, "could not save category", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, loggedIn := middleware.IsAuthenticated(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		HandleError(w, "bad request", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {
		HandleError(w, "bad request", http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := models.InsertComment(models.Comments{
		UserId:  userID,
		PostId:  postID,
		Content: content,
	}); err != nil {
		HandleError(w, "could not create comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ReactPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, loggedIn := middleware.IsAuthenticated(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		HandleError(w, "bad request", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {
		HandleError(w, "bad request", http.StatusBadRequest)
		return
	}

	reactionType := r.FormValue("type")
	if reactionType != "like" && reactionType != "dislike" {
		HandleError(w, "bad request", http.StatusBadRequest)
		return
	}

	currentReaction, err := models.CheckReactionByUser(userID, postID, nil)
	if err != nil {
		HandleError(w, "could not update reaction", http.StatusInternalServerError)
		return
	}

	if currentReaction == reactionType {
		if err := models.DeleteReaction(userID, postID, nil); err != nil {
			HandleError(w, "could not remove reaction", http.StatusInternalServerError)
			return
		}
	} else {
		if currentReaction != "" {
			if err := models.DeleteReaction(userID, postID, nil); err != nil {
				HandleError(w, "could not replace reaction", http.StatusInternalServerError)
				return
			}
		}

		_, err = models.InsertReaction(models.Reaction{
			UserID: userID,
			PostID: postID,
			Type:   reactionType,
		})
		if err != nil {
			HandleError(w, "could not save reaction", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func loadPosts(r *http.Request, userID int, loggedIn bool) ([]models.Post, int, string, error) {
	selectedCategoryID := 0
	selectedView := strings.TrimSpace(r.URL.Query().Get("view"))
	if selectedView != "mine" && selectedView != "liked" {
		selectedView = ""
	}
	if !loggedIn {
		selectedView = ""
	}

	if categoryParam := strings.TrimSpace(r.URL.Query().Get("category")); categoryParam != "" {
		id, err := strconv.Atoi(categoryParam)
		if err == nil && id > 0 {
			selectedCategoryID = id
		}
	}

	var (
		posts []models.Post
		err   error
	)

	if selectedCategoryID > 0 {
		posts, err = models.GetPostsByCategory(selectedCategoryID)
	} else {
		posts, err = models.GetAllPosts()
	}
	if err != nil {
		return nil, 0, "", err
	}

	enrichedPosts := make([]models.Post, 0, len(posts))
	for _, post := range posts {
		post.Likes, _ = models.CountLikeDislikeByPost(post.IdPost, "like")
		post.Dislikes, _ = models.CountLikeDislikeByPost(post.IdPost, "dislike")
		post.Comments, _ = models.GetCommentsByPost(post.IdPost)

		if loggedIn {
			reaction, _ := models.GetReactionByUser(userID, post.IdPost)
			post.UserReaction = reaction

			switch selectedView {
			case "mine":
				if post.UserId != userID {
					continue
				}
			case "liked":
				if reaction != "like" {
					continue
				}
			}
		}

		enrichedPosts = append(enrichedPosts, post)
	}

	return enrichedPosts, selectedCategoryID, selectedView, nil
}
