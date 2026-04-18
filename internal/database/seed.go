package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedData() {
	seedCategories()
	seedAdminUser()
}

func seedCategories() {
	categories := []string{
		"General",
		"Announcements",
		"Tech",
		"Gaming",
		"Help",
	}

	for _, name := range categories {
		var count int
		err := DB.QueryRow("SELECT COUNT(*) FROM category WHERE name = ?", name).Scan(&count)
		if err != nil {
			log.Println("seed category check:", err)
			continue
		}

		if count > 0 {
			continue
		}

		if _, err := DB.Exec("INSERT INTO category (name) VALUES (?)", name); err != nil {
			log.Println("seed category insert:", err)
		}
	}
}

func seedAdminUser() {
	var id int
	err := DB.QueryRow("SELECT id FROM users WHERE email = ?", "admin@kuu.ku").Scan(&id)
	if err == nil {
		return
	}
	if err != sql.ErrNoRows {
		log.Println("seed admin check:", err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Println("seed admin hash:", err)
		return
	}

	_, err = DB.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		"admin",
		"admin@kuu.ku",
		string(hash),
	)
	if err != nil {
		log.Println("seed admin insert:", err)
	}
}
