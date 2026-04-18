package models

import (
	"log"
	"time"

	"forum/database"
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
              ON posts.user_id = users.id`

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

func GetPostById(id int) (Post, error) {
	post := Post{}
	query := "SELECT id, title, content, user_id, image, created_at  FROM posts WHERE id = ?"
	err := database.DB.QueryRow(query, id).Scan(&post.IdPost, &post.Title, &post.Content, &post.UserId, &post.Image, &post.CreatedAt)
	if err != nil {
		return Post{}, err
	}
	return post, nil
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

func FetchComment(postId int) ([]Comments, error) {
	comments := []Comments{}
	query := `SELECT comments.id, comments.user_id, comments.post_id,comments.content,comments.created_at,users.username
	FROM comments INNER JOIN users ON comments.user_id=users.id
			WHERE comments.post_id=? `

	rows, err := database.DB.Query(query, postId)
	if err != nil {
		log.Println("FetchComment error:", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		c := Comments{}
		err := rows.Scan(&c.IdComment, &c.UserId, &c.PostId, &c.Content, &c.CreatedAt, &c.Username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
