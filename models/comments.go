package models

import "time"

type Comments struct {
	id_comment int 
	user_id int
	post_id int
	content string
	created_at time.Time
}