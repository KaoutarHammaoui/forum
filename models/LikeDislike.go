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
func DeleteReaction(userId, postId int) error {
	query := "DELETE FROM likes_dislikes WHERE user_id = ? AND post_id = ?"
	_, err := database.DB.Exec(query, userId, postId)
	if err != nil {
		return err
	}
	return nil
}

func CheckReactionByUser(userId, postId int) (bool, error) {
	exist := 0
	query := "SELECT COUNT(*) FROM likes_dislikes WHERE user_id = ? AND post_id = ?"
	err := database.DB.QueryRow(query, userId, postId).Scan(&exist)
	if err != nil {
		return false, err
	}
	if exist == 0 {
		return false, nil
	}
	return true, nil
}

func GetReactionByUser(userId, postId int) (string, error) {
	reactionType := ""
	query := "SELECT type FROM likes_dislikes WHERE user_id = ? AND post_id = ? LIMIT 1"
	err := database.DB.QueryRow(query, userId, postId).Scan(&reactionType)
	if err != nil {
		return "", nil // pas de réaction trouvée
	}
	return reactionType, nil
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
