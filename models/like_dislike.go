package models

import "time"

type LikeDislike struct {
	id_lD int
	user_id int
	post_id int
	content string
	created_at time.Time
	type_ld string
}