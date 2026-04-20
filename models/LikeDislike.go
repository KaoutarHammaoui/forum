package models

import (
	"forum/database"
)

type Reaction struct {
	ID        int
	UserID    int
	PostID    int
	CommentID *int
	Type      string
}

func InsertReaction(reaction Reaction) (int64, error) {
	query := "INSERT INTO likes_dislikes (user_id, post_id, comment_id, type) VALUES (?, ?, ?, ?)"
	result, err := database.DB.Exec(query, reaction.UserID, reaction.PostID, reaction.CommentID, reaction.Type)
	if err != nil {
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil

}

func DeleteReaction(userId int, postID int) error {
	query := `
	DELETE FROM likes_dislikes 
	WHERE user_id=? AND post_id=? AND comment_id IS NULL
	`
	_, err := database.DB.Exec(query, userId, postID)
	return err
}

func CountLikeDislikeByPost(postId int, Type string) (int, error) {
	count := 0
	query := "SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND type = ?"
	err := database.DB.QueryRow(query, postId, Type).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetReactionByUser(userId, postId int) (string, error) {
	reactionType := ""
	query := `
	SELECT type FROM likes_dislikes 
	WHERE user_id = ? AND post_id = ? AND comment_id IS NULL
	LIMIT 1
	`
	err := database.DB.QueryRow(query, userId, postId).Scan(&reactionType)
	if err != nil {
		return "", nil // pas de réaction trouvée
	}
	return reactionType, nil
}
