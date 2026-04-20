package handlers

import (
	"forum/middleware"
	"forum/models"
	"net/http"
	"strconv"
	"strings"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", 405)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIdKey).(int)
	if !ok {
		HandleError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		HandleError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Redirect(w, r, "/homeUser", 302)
		return
	}

	comment := models.Comments{
		UserId:  userID,
		PostId:  postID,
		Content: content,
	}

	err = models.InsertComment(comment)
	if err != nil {
		HandleError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/homeUser", 302)
}
