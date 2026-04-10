package models

import (
	"forum/database"
	"time"
)

type LikeDislike struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

func AddReaction(userID, postID, commentID int, reactionType string) error {
	_, _ = database.DB.Exec(`
		DELETE FROM likes_dislikes 
		WHERE user_id = ? AND (post_id = ? OR comment_id = ?)`,
		userID, postID, commentID)

	_, err := database.DB.Exec(`
		INSERT INTO likes_dislikes (user_id, post_id, comment_id, type)
		VALUES (?, ?, ?, ?)`,
		userID, postID, commentID, reactionType)

	return err
}

func CountReactions(postID, commentID int, reactionType string) (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM likes_dislikes 
	          WHERE type = ? AND (post_id = ? OR comment_id = ?)`

	err := database.DB.QueryRow(query, reactionType, postID, commentID).Scan(&count)

	return count, err
}
