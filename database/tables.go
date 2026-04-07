package database

import "log"

func TableCreation() {
	queries := []string{
	`CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY,
  			username VARCHAR UNIQUE NOT NULL,
  			email VARCHAR UNIQUE NOT NULL,
  			password VARCHAR NOT NULL,
  			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`,

	`CREATE TABLE IF NOT EXISTS posts(
			id INTEGER PRIMARY KEY,
			TITLE VARCHAR,
			CONTENT TEXT,
			user_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
	)`,

	`CREATE TABLE IF NOT EXISTS category(
		id INTEGER PRIMARY KEY,
		name VARCHAR NOT NULL
	)`,

	`CREATE TABLE IF NOT EXISTS post_category(
		post_id INTEGER ,
		category_id INTEGER,
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY (category_id) REFERENCES category(id)
	)`,

	`CREATE TABLE IF NOT EXISTS comments(
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		content VARCHAR,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(post_id) REFERENCES posts(id)
	)`,

	`CREATE TABLE IF NOT EXISTS likes_dislikes(		id INTEGER PRIMARY KEY,
		id INTEGER PRIMARY KEY,
		user_id INTEGER,
		post_id INTEGER,
		comment_id INTEGER,
		type 	TEXT NOT NULL CHECK(type IN ('like','dislike')),
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY (comment_id) REFERENCES comments(id)
	)`,

	`CREATE TABLE IF NOT EXISTS session(
		id INTEGER PRIMARY KEY,
		user_id INTEGER,
		token VARCHAR UNIQUE NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`,

	}
	for _,query:=range queries{
		prep,err:=DB.Prepare(query) //to avoid sql injections
		if err!=nil{
			log.Fatal(err)
		}
		_,err=prep.Exec()

		if err !=nil{
			log.Fatal(err)
		}
	}
}
