package models

import "time"

//On utilise le Model  afin de transformer une table a un objet manipuable . 
//Creation  d Model liée a la table Users : 
	
type Users struct {
	id_user			int
	username 	string 
	email 		string
	password 	string
	created_at  time.Time
}


// CRUD  => Create + Find  ; Authentification 


func InsertUser(){
	rqst := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	
}