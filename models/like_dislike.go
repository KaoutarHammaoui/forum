package models

import "forum/database"

type LikeDislike struct {
	ID        int
	UserID    int
	PostID    int
	CommentID *int
	Type      string 
}

func SetReaction(userID, postID int, commentID *int, reactionType string) error {
	_, err := database.DB.Exec(`
		INSERT INTO likes_dislikes (user_id, post_id, comment_id, type)
		VALUES (?, ?, ?, ?)`,
		userID, postID, commentID, reactionType,
	)
	return err
}