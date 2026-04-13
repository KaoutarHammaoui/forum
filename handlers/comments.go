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
		http.Error(w, "Method not allowed", 405)
		return
	}

	userID := r.Context().Value(middleware.UserIdKey).(int)
	postID, _ := strconv.Atoi(r.FormValue("post_id"))
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
