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
			title VARCHAR,
			content TEXT,
			user_id INTEGER NOT NULL,
			image TEXT, 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`,

		`CREATE TABLE IF NOT EXISTS category(
			id INTEGER PRIMARY KEY,
			name VARCHAR NOT NULL
	)`,

		`CREATE TABLE IF NOT EXISTS post_category(
			post_id INTEGER ,
			category_id INTEGER,
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
	)`,

		`CREATE TABLE IF NOT EXISTS comments(
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		content VARCHAR,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
	)`,

		`CREATE TABLE IF NOT EXISTS likes_dislikes(		
		id INTEGER PRIMARY KEY,
		user_id INTEGER,
		post_id INTEGER,
		comment_id INTEGER,
		type 	TEXT NOT NULL CHECK(type IN ('like','dislike')),
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
	)`,

		`CREATE TABLE IF NOT EXISTS session(
		id INTEGER PRIMARY KEY,
		user_id INTEGER,
		token VARCHAR UNIQUE NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`,
	}
	for _, query := range queries {
		prep, err := DB.Prepare(query) // to avoid sql injections
		if err != nil {
			log.Fatal(err)
		}
		_, err = prep.Exec()
		if err != nil {
			log.Fatal(err)
		}
	}
	seedQueries := []string{
		// USERS
		`INSERT OR IGNORE INTO users (id, username, email, password) VALUES
		(1, 'oumaima', 'oumaima@mail.com', '1234'),
		(2, 'kaoutar', 'kaoutar@mail.com', '1234'),
		(3, 'admin', 'admin@mail.com', '1234')`,

		// CATEGORY
		`INSERT OR IGNORE INTO category (id, name) VALUES
		(1, 'Technologie'),
		(2, 'Science'),
		(3, 'Art'),
		(4, 'Music'),
		(5, 'Animal')`,
		
	}

	// INSERT DATA
	for _, query := range seedQueries {
		_, err := DB.Exec(query)
		if err != nil {
			log.Println("Seed error:", err)
		}
	}

	log.Println("Seed data inserted successfully")
}
