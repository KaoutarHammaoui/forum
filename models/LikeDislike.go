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
	query := "INSERT INTO likes_dislikes (user_id, post_id, comment_id, type) VALUES (?, ?, ?, ?)`"
	result, err := database.DB.Exec(query, reaction.UserID, reaction.PostID, reaction.CommentID, reaction.Type)
	if err != nil {
		return 0, err
	}
	lastId,_ := result.LastInsertId()
	return lastId, nil

}
func DeleteReaction() error{
	return nil
}

func CheckReactionByUser(userId int) error {
	return nil
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
