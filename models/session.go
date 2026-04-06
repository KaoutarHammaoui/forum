package models

import "time"

type Session struct {
	id_session int
	user_id int
	token string
	expires_at time.Time
	
}