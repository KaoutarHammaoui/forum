package handlers

import (
	"forum/middleware"
	"forum/models"
	"net/http"
	"strconv"
)

func ReactPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", 405)
		return
	}

	userID := r.Context().Value(middleware.UserIdKey).(int)
	postID, _ := strconv.Atoi(r.FormValue("post_id"))
	reactionType := r.FormValue("type")

	// Vérifier si une réaction existe déjà
	existingType, err := models.GetReactionByUser(userID, postID)
	if err != nil {
		HandleError(w, "Internal Server Error", 500)
		return
	}

	if existingType == reactionType {
		// Même réaction → on supprime (toggle off)
		models.DeleteReaction(userID, postID)
	} else {
		// Pas de réaction ou réaction différente → supprimer l'ancienne et insérer
		models.DeleteReaction(userID, postID)
		reaction := models.Reaction{
			UserID:    userID,
			PostID:    postID,
			Type:      reactionType,
			CommentID: nil,
		}
		_, err := models.InsertReaction(reaction)
		if err != nil {
			HandleError(w, "Internal Server Error", 500)
			return
		}
	}

	http.Redirect(w, r, r.FormValue("redirect"), 302)
}
