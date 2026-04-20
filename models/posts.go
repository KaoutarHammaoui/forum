package models

import (
	"forum/database"
	"strings"
	"time"
)

type Post struct {
	IdPost    int
	Title     string
	Content   string
	UserId    int
	Image     string
	UserName  string
	Comments  []Comments
	Likes     int
	Dislikes  int
	CreatedAt time.Time
}

func InsertPost(post Post) (int64, error) {
	query := "INSERT INTO posts (title, content, user_id, image) VALUES (?, ?, ?, ?)"
	result, err := database.DB.Exec(query, post.Title, post.Content, post.UserId, post.Image)
	if err != nil {
		return 0, err
	}
	lastPostId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastPostId, nil
}

func GetAllPosts() ([]Post, error) {
	posts := []Post{}
	query := `SELECT posts.id, posts.title, posts.content, posts.user_id, posts.image, posts.created_at, users.username
              FROM posts
              INNER JOIN users
              ON posts.user_id = users.id ORDER BY posts.created_at DESC `

	lignes, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer lignes.Close()

	for lignes.Next() {
		post := Post{}
		err := lignes.Scan(&post.IdPost, &post.Title, &post.Content, &post.UserId, &post.Image, &post.CreatedAt, &post.UserName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := lignes.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostsByCategory(idcat int) ([]Post, error) {
	posts := []Post{}
	query := `SELECT p.id, p.title, p.content, p.user_id, p.image, p.created_at, u.username
			  FROM posts p
			  INNER JOIN post_category pc ON p.id = pc.post_id
			  INNER JOIN users u ON p.user_id = u.id
			  WHERE pc.category_id = ?`

	lignes, err := database.DB.Query(query, idcat)
	if err != nil {
		return nil, err
	}
	defer lignes.Close()

	for lignes.Next() {
		post := Post{}
		err := lignes.Scan(&post.IdPost, &post.Title, &post.Content, &post.UserId, &post.Image, &post.CreatedAt, &post.UserName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := lignes.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetFilteredPosts(userID int, selectedCats []string, filterLikes, filterMyPosts bool) ([]Post, error) {
	var posts []Post

	query := `SELECT DISTINCT p.id, p.title, p.content, p.user_id, p.image, p.created_at, u.username
              FROM posts p
              INNER JOIN users u ON p.user_id = u.id`

	var conditions []string
	var args []interface{}

	if filterLikes {
		query += " INNER JOIN likes_dislikes ld ON p.id = ld.post_id"
		conditions = append(conditions, "ld.user_id = ? AND ld.type = ?")
		args = append(args, userID, "like")
	}

	filteredCats := make([]string, 0, len(selectedCats))
	for _, cat := range selectedCats {
		if cat != "all" {
			filteredCats = append(filteredCats, cat)
		}
	}
	if len(filteredCats) > 0 {
		query += " INNER JOIN post_category pc ON p.id = pc.post_id"

		placeholders := make([]string, 0, len(filteredCats))
		for _, cat := range filteredCats {
			placeholders = append(placeholders, "?")
			args = append(args, cat)
		}
		conditions = append(conditions, "pc.category_id IN ("+strings.Join(placeholders, ",")+")")
	}

	if filterMyPosts {
		conditions = append(conditions, "p.user_id = ?")
		args = append(args, userID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY p.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.IdPost, &p.Title, &p.Content, &p.UserId, &p.Image, &p.CreatedAt, &p.UserName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
