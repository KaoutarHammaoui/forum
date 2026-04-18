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


	_,err := models.InsertReaction(reaction)
	if err != nil {
	
		return
	}

	http.Redirect(w, r, "/homeUser", 302)
}
