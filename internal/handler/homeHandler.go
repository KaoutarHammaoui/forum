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
	Title               string
	IsLoggedIn          bool
	UserID              int
	Categories          []models.Category
	Posts               []models.Post
	SelectedCategoryIDs []int
	SelectedCategoryMap map[int]bool
	SelectedView        string
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

	posts, selectedCategoryIDs, selectedView, err := loadPosts(r, userID, loggedIn)
	if err != nil {
		HandleError(w, "could not load posts", http.StatusInternalServerError)
		return
	}

	selectedCategoryMap := make(map[int]bool, len(selectedCategoryIDs))
	for _, categoryID := range selectedCategoryIDs {
		selectedCategoryMap[categoryID] = true
	}

	data := Data{
		Title:               "Home",
		IsLoggedIn:          loggedIn,
		UserID:              userID,
		Categories:          categories,
		Posts:               posts,
		SelectedCategoryIDs: selectedCategoryIDs,
		SelectedCategoryMap: selectedCategoryMap,
		SelectedView:        selectedView,
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
	categoryValues := r.Form["category_id"]

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

	seenCategories := map[int]bool{}
	for _, rawCategoryID := range categoryValues {
		categoryID, err := strconv.Atoi(rawCategoryID)
		if err != nil || categoryID <= 0 || seenCategories[categoryID] {
			continue
		}

		seenCategories[categoryID] = true
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

	var commentID *int
	if rawCommentID := strings.TrimSpace(r.FormValue("comment_id")); rawCommentID != "" {
		parsedCommentID, err := strconv.Atoi(rawCommentID)
		if err != nil || parsedCommentID <= 0 {
			HandleError(w, "bad request", http.StatusBadRequest)
			return
		}
		commentID = &parsedCommentID
	}

	reactionType := r.FormValue("type")
	if reactionType != "like" && reactionType != "dislike" {
		HandleError(w, "bad request", http.StatusBadRequest)
		return
	}

	currentReaction, err := models.CheckReactionByUser(userID, postID, commentID)
	if err != nil {
		HandleError(w, "could not update reaction", http.StatusInternalServerError)
		return
	}

	if currentReaction == reactionType {
		if err := models.DeleteReaction(userID, postID, commentID); err != nil {
			HandleError(w, "could not remove reaction", http.StatusInternalServerError)
			return
		}
	} else {
		if currentReaction != "" {
			if err := models.DeleteReaction(userID, postID, commentID); err != nil {
				HandleError(w, "could not replace reaction", http.StatusInternalServerError)
				return
			}
		}

		_, err = models.InsertReaction(models.Reaction{
			UserID:    userID,
			PostID:    postID,
			CommentID: commentID,
			Type:      reactionType,
		})
		if err != nil {
			HandleError(w, "could not save reaction", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func loadPosts(r *http.Request, userID int, loggedIn bool) ([]models.Post, []int, string, error) {
	selectedCategoryIDs := []int{}
	selectedView := strings.TrimSpace(r.URL.Query().Get("view"))
	if selectedView != "mine" && selectedView != "liked" {
		selectedView = ""
	}
	if !loggedIn {
		selectedView = ""
	}

	seenCategories := map[int]bool{}
	for _, categoryParam := range r.URL.Query()["category"] {
		id, err := strconv.Atoi(strings.TrimSpace(categoryParam))
		if err == nil && id > 0 && !seenCategories[id] {
			selectedCategoryIDs = append(selectedCategoryIDs, id)
			seenCategories[id] = true
		}
	}

	var (
		posts []models.Post
		err   error
	)

	if len(selectedCategoryIDs) > 0 {
		posts, err = models.GetPostsByCategories(selectedCategoryIDs)
	} else {
		posts, err = models.GetAllPosts()
	}
	if err != nil {
		return nil, nil, "", err
	}

	enrichedPosts := make([]models.Post, 0, len(posts))
	for _, post := range posts {
		post.Likes, _ = models.CountLikeDislikeByPost(post.IdPost, "like")
		post.Dislikes, _ = models.CountLikeDislikeByPost(post.IdPost, "dislike")
		post.Comments, _ = models.GetCommentsByPost(post.IdPost)

		if loggedIn {
			reaction, _ := models.GetReactionByUser(userID, post.IdPost)
			post.UserReaction = reaction
			for i := range post.Comments {
				commentReaction, _ := models.CheckReactionByUser(userID, post.IdPost, &post.Comments[i].IdComment)
				post.Comments[i].UserReaction = commentReaction
			}

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

	return enrichedPosts, selectedCategoryIDs, selectedView, nil
}
