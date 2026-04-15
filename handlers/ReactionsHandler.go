package handlers

import (
	"net/http"
	"strconv"

	"forum/middleware"
	"forum/models"
)

func ReactPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", 405)
		return
	}

	userID := r.Context().Value(middleware.UserIdKey).(int)

	postID, _ := strconv.Atoi(r.FormValue("post_id"))
	reactionType := r.FormValue("type") // like / dislike

	err := models.SetReaction(userID, postID, nil, reactionType)
	if err != nil {
		HandleError(w, "DB error", 500)
		return
	}

	http.Redirect(w, r, "/homeUser", 302)
}