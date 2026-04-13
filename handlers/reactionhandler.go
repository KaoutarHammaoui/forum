package handlers

import (
	"net/http"
	"strconv"

	"forum/middleware"
	"forum/models"
)

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(middleware.UserIdKey).(int)

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil {
		commentID = 0
	}
	reactionType := r.FormValue("type")

	if reactionType != "like" && reactionType != "dislike" {
		http.Error(w, "Invalid reaction", http.StatusBadRequest)
		return
	}

	err = models.AddReaction(userID, postID, commentID, reactionType)
	if err != nil {
		http.Error(w, "Error reaction", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
