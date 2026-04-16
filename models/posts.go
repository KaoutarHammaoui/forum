package models

import (
	"strings"
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
	CreatedAt time.Time
}

// A changer :!!!
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

func GetFilteredPosts(userID int, selectedCats []string, filterLikes, filterMyPosts bool) ([]Post, error) {
    var posts []Post
    
    // Base Query
    query := `SELECT DISTINCT p.id, p.title, p.content, p.user_id, p.image, p.created_at, u.username
              FROM posts p
              INNER JOIN users u ON p.user_id = u.id`
    
    var conditions []string
    var args []interface{}

    // 1. Filter by Likes
    if filterLikes {
        query += " INNER JOIN likes_dislikes ld ON p.id = ld.post_id"
        conditions = append(conditions, "ld.user_id = ? AND ld.type = 'like' AND ld.type = 'dislike'")
        args = append(args, userID)
    }

    // 2. Filter by Categories
    if len(selectedCats) > 0 {
        query += " INNER JOIN post_category pc ON p.id = pc.post_id"
        
        // Create placeholders (?, ?, ?) for the categories
        placeholders := []string{}
        for _, cat := range selectedCats {
            if cat != "all" {
                placeholders = append(placeholders, "?")
                args = append(args, cat)
            }
        }
        if len(placeholders) > 0 {
            conditions = append(conditions, "pc.category_id IN (" + strings.Join(placeholders, ",") + ")")
        }
    }

    // 3. Filter by User's Own Posts
    if filterMyPosts {
        conditions = append(conditions, "p.user_id = ?")
        args = append(args, userID)
    }

    // Combine all WHERE conditions
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
    return posts, nil
}