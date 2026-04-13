package models

import "time"

type Comments struct {
	idComment int
	userId    int
	postId    int
	content   string
	createdAt time.Time
}
