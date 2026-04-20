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

	userID := r.Context().Value(middleware.UserIdKey).(int)
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

	models.InsertComment(comment)
	http.Redirect(w, r, "/homeUser", 302)
}

