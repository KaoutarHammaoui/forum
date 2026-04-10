package handlers

import (
	"net/http"
	"strconv"

	"forum/models"
)

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	postID, _ := strconv.Atoi(r.FormValue("post_id"))
	commentID, _ := strconv.Atoi(r.FormValue("comment_id"))
	reactionType := r.FormValue("type") // like or dislike

	if reactionType != "like" && reactionType != "dislike" {
		http.Error(w, "Invalid reaction", http.StatusBadRequest)
		return
	}

	err := models.AddReaction(userID, postID, commentID, reactionType)
	if err != nil {
		http.Error(w, "Error reaction", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
