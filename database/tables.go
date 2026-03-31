package database

func TableCreation() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users(
		 id INTEGER PRIMARY KEY,
  username VARCHAR UNIQUE,
  email VARCHAR UNIQUE,
  password VARCHAR,
  created_at TIMESTAMP
	)`,
	}
}
