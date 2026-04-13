package handlers

import (
	"net/http"
	"strconv"

	"forum/middleware"
	"forum/models"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(middleware.UserIdKey).(int)

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	content := r.FormValue("content")

	if content == "" {
		http.Error(w, "Empty comment", http.StatusBadRequest)
		return
	}

	err = models.CreateComment(userID, postID, content)
	if err != nil {
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
