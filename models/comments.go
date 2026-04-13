package models

import (
        "forum/database"
        "time"
)

type Comments struct {
        IdComment int
        UserId    int
        PostId    int
        Content   string
        CreatedAt time.Time
        Username  string
}

func InsertComment(c Comments) error {
        _, err := database.DB.Exec(
                "INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)",
                c.UserId, c.PostId, c.Content,
        )
        return err
}

func GetCommentsByPost(postID int) ([]Comments, error) {
        rows, err := database.DB.Query(`
                SELECT c.id, c.user_id, c.post_id, c.content, c.created_at, u.username
                FROM comments c
                JOIN users u ON u.id = c.user_id
                WHERE c.post_id = ?`, postID)

        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var comments []Comments

        for rows.Next() {
                var c Comments

                err := rows.Scan(
                        &c.IdComment,
                        &c.UserId,
                        &c.PostId,
                        &c.Content,
                        &c.CreatedAt,
                        &c.Username,
                )
                if err != nil {
                        return nil, err
                }

                comments = append(comments, c)
        }

        if err := rows.Err(); err != nil {
                return nil, err
        }

        return comments, nil
}
