package models

import "time"

type Posts  struct {
	id_post int 
	title string 
	content string
	user_id int 
	created_at time.Time
	FK int
}