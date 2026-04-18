package handlers

import (
	"net/http"
	"strconv"

	"forum/middleware"
	"forum/models"
)

func ReactPost(w http.ResponseWriter, r *http.Request) {
		reaction := models.Reaction{}
	if r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", 405)
		return
	}

	userID := r.Context().Value(middleware.UserIdKey).(int)
	postID, _ := strconv.Atoi(r.FormValue("post_id"))
	reactionType := r.FormValue("type")

	reaction.UserID = userID
	reaction.PostID = postID
	reaction.Type = reactionType
	reaction.CommentID = nil

	// Vérifier si une réaction existe déjà
	existingType, err := models.GetReactionByUser(userID, postID)
	if err != nil {
		HandleError(w, "Internal Server Error", 500)
		return
	}

	if existingType == reactionType {
		// Même réaction → on supprime (toggle off)
		models.DeleteReaction(userID, postID,reaction.CommentID)
	} else {
		// Pas de réaction ou réaction différente → supprimer l'ancienne et insérer
		models.DeleteReaction(userID, postID,reaction.CommentID)
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

func ReactComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, "Method not allowed", 405)
		return
	}
	UserID := r.Context().Value(middleware.UserIdKey).(int)
	postId := r.FormValue("postID")
	commentId := r.FormValue("commentID")
	reactionType := r.FormValue("type")

	postID, err := strconv.Atoi(postId)
	if err != nil {
		HandleError(w, "error in converting post id to int", http.StatusBadRequest)
		return
	}
	commentID, er := strconv.Atoi(commentId)
	if er != nil {
		HandleError(w, "invalid comment id", http.StatusBadRequest)
		return
	}
	commentIDPtr := &commentID

	reactionChecking, err := models.CheckReactionByUser(UserID, postID, commentIDPtr)
	if err != nil {
		HandleError(w, "error during checking reaction", 500)
		return
	}
	if reactionChecking == "" {
		_, err := models.InsertReaction(models.Reaction{
			UserID:    UserID,
			PostID:    postID,
			CommentID: commentIDPtr,
			Type:      reactionType,
		})
		if err != nil {
			HandleError(w, "delete error", 500)
			return
		}

	} else if reactionChecking == reactionType {
		err := models.DeleteReaction(UserID, postID, &commentID)
		if err != nil {
			HandleError(w, "delete error", 500)
			return
		}
	} else {
		err := models.DeleteReaction(UserID, postID, &commentID)
		if err != nil {
			HandleError(w, "delete error", 500)
			return
		}
		_, er := models.InsertReaction(models.Reaction{
			UserID:    UserID,
			PostID:    postID,
			CommentID: commentIDPtr,
			Type:      reactionType,
		})
		if er != nil {
			HandleError(w, "delete error", 500)
			return
		}
	}
	http.Redirect(w, r, "/homeUser", 302)

}
