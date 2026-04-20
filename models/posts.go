package models

import (
	"forum/database"
	"sort"
	"time"
)

type Post struct {
	IdPost    int
	Title     string
	Content   string
	UserId    int
	UserName  string
	Comments  []Comments
	Likes     int
	Dislikes  int
	CreatedAt time.Time
}

func InsertPost(post Post) (int64, error) {
	query := "INSERT INTO posts (title, content, user_id) VALUES (?, ?, ?)"
	result, err := database.DB.Exec(query, post.Title, post.Content, post.UserId)
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
	query := `SELECT posts.id, posts.title, posts.content, posts.user_id, posts.created_at, users.username
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
		err := lignes.Scan(&post.IdPost, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, &post.UserName)
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
	query := `SELECT p.id, p.title, p.content, p.user_id, p.created_at, u.username
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
		err := lignes.Scan(&post.IdPost, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, &post.UserName)
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
	seen := map[int]bool{}
	allPosts := []Post{}

	addPosts := func(posts []Post) {
		for _, p := range posts {
			if !seen[p.IdPost] {
				seen[p.IdPost] = true
				allPosts = append(allPosts, p)
			}
		}
	}

	// Filtre catégories
	for _, cat := range selectedCats {
		if cat == "all" {
			continue
		}
		rows, err := database.DB.Query(`
			SELECT DISTINCT p.id, p.title, p.content, p.user_id, p.created_at, u.username
			FROM posts p
			INNER JOIN users u ON p.user_id = u.id
			INNER JOIN post_category pc ON p.id = pc.post_id
			WHERE pc.category_id = ?
			ORDER BY p.created_at DESC`, cat)
		if err != nil {
			return nil, err
		}
		var posts []Post
		for rows.Next() {
			var p Post
			if err := rows.Scan(&p.IdPost, &p.Title, &p.Content, &p.UserId, &p.CreatedAt, &p.UserName); err != nil {
				rows.Close()
				return nil, err
			}
			posts = append(posts, p)
		}
		rows.Close()
		addPosts(posts)
	}

	// Filtre liked posts
	if filterLikes {
		rows, err := database.DB.Query(`
			SELECT DISTINCT p.id, p.title, p.content, p.user_id, p.created_at, u.username
			FROM posts p
			INNER JOIN users u ON p.user_id = u.id
			INNER JOIN likes_dislikes ld ON p.id = ld.post_id
			WHERE ld.user_id = ? AND ld.type = ?
			ORDER BY p.created_at DESC`, userID, "like")
		if err != nil {
			return nil, err
		}
		var posts []Post
		for rows.Next() {
			var p Post
			if err := rows.Scan(&p.IdPost, &p.Title, &p.Content, &p.UserId, &p.CreatedAt, &p.UserName); err != nil {
				rows.Close()
				return nil, err
			}
			posts = append(posts, p)
		}
		rows.Close()
		addPosts(posts)
	}

	// Filtre my posts
	if filterMyPosts {
		rows, err := database.DB.Query(`
			SELECT p.id, p.title, p.content, p.user_id, p.created_at, u.username
			FROM posts p
			INNER JOIN users u ON p.user_id = u.id
			WHERE p.user_id = ?
			ORDER BY p.created_at DESC`, userID)
		if err != nil {
			return nil, err
		}
		var posts []Post
		for rows.Next() {
			var p Post
			if err := rows.Scan(&p.IdPost, &p.Title, &p.Content, &p.UserId, &p.CreatedAt, &p.UserName); err != nil {
				rows.Close()
				return nil, err
			}
			posts = append(posts, p)
		}
		rows.Close()
		addPosts(posts)
	}

	sort.Slice(allPosts, func(i, j int) bool {
		return allPosts[i].CreatedAt.After(allPosts[j].CreatedAt)
	})

	return allPosts, nil
}
