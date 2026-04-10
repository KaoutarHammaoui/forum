package models

import (
	"forum/database"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateComment(userID, postID int, content string) error {
	query := `INSERT INTO comments (user_id, post_id, content, created_at)
	          VALUES (?, ?, ?, ?)`
	_, err := database.DB.Exec(query, userID, postID, content, time.Now())
	return err
}

func GetCommentsByPost(postID int) ([]Comment, error) {
	rows, err := database.DB.Query(`SELECT id, user_id, post_id, content, created_at 
		FROM comments WHERE post_id = ? ORDER BY created_at DESC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.ID, &c.UserID, &c.PostID, &c.Content, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
