package models

import "time"

type Posts struct {
	idPost    int
	title     string
	content   string
	userId    int
	createdAt time.Time
}

/*
	CreatePost
	GetAllPosts
	GetPostByID
	GetPostsByUserID
	GetPostsByCategory
	GetLikedPostsByUserID
*/
