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
		HandleError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIdKey).(int)
	if !ok {
		HandleError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		HandleError(w, "Bad Request", http.StatusBadRequest)
		return
	}

	reactionType := r.FormValue("type")
	if reactionType != "like" && reactionType != "dislike" {
		HandleError(w, "Bad Request", http.StatusBadRequest)
		return
	}

	reaction.UserID = userID
	reaction.PostID = postID
	reaction.Type = reactionType
	reaction.CommentID = nil

	existingType, err := models.GetReactionByUser(userID, postID)
	if err != nil {
		HandleError(w, "Internal Server Error", 500)
		return
	}

	if existingType == reactionType {
		err := models.DeleteReaction(userID, postID, reaction.CommentID)
		if err != nil {
			HandleError(w, "Internal Server Error", 500)
			return
		}
	} else {
		err = models.DeleteReaction(userID, postID, reaction.CommentID)
		if err != nil {
			HandleError(w, "Internal Server Error", 500)
			return
		}
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

	UserID, ok := r.Context().Value(middleware.UserIdKey).(int)
	if !ok {
		HandleError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("postID"))
	if err != nil {
		HandleError(w, "error in converting post id to int", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.Atoi(r.FormValue("commentID"))
	if err != nil {
		HandleError(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	reactionType := r.FormValue("type")
	if reactionType != "like" && reactionType != "dislike" {
		HandleError(w, "Bad Request", http.StatusBadRequest)
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
			HandleError(w, "Internal Server Error", 500)
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
		_, err = models.InsertReaction(models.Reaction{
			UserID:    UserID,
			PostID:    postID,
			CommentID: commentIDPtr,
			Type:      reactionType,
		})
		if err != nil {
			HandleError(w, "delete error", 500)
			return
		}
	}

	http.Redirect(w, r, "/homeUser", 302)
}
