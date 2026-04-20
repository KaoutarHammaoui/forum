package handlers

import (
	"fmt"

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
	userID, ok := r.Context().Value(middleware.UserIdKey).(int)
	if !ok {
		HandleError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		HandleError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	content := strings.TrimSpace(r.FormValue("Commentcontent"))
	if content == "" {
		http.SetCookie(w, &http.Cookie{
			Name: "comment_error",
			Value: fmt.Sprintf("%d", postID),
			Path: "/",
		})

		http.Redirect(w, r, "/homeUser", http.StatusSeeOther)
		return
	}

	comment := models.Comments{
		UserId: userID,
		PostId: postID,
		Content: content,
	}
	err = models.InsertComment(comment)
	if err != nil {
		HandleError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/homeUser", 302)

}
