package models

import "time"

type LikeDislike struct {
	idlD int
	userId int
	postId int
	content string
	createdAt time.Time
	typeLd string
}