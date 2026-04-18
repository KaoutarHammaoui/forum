package models

import (
	"forum/internal/database"
	"time"
)

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

func InsertUser(user User) (int64, error) {
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	result, err := database.DB.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

func ExistsInColumn(column, value string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users  WHERE " + column + " = ?"
	err := database.DB.QueryRow(query, value).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUserByEmail(email string) (User, error) {
	payload := User{}
	query := "SELECT id, username, email, password FROM users  WHERE email = ?"
	err := database.DB.QueryRow(query, email).Scan(&payload.ID, &payload.Username, &payload.Email, &payload.Password)
	if err != nil {
		return User{}, err
	}
	return payload, nil
}


func GetUserByID(id int) (User, error) {
	payload := User{}
	query := "SELECT id, username, email, password, created_at  FROM users  WHERE id = ?"
	err := database.DB.QueryRow(query, id).Scan(&payload.ID, &payload.Username, &payload.Email, &payload.Password, &payload.CreatedAt)
	if err != nil {
		return User{}, err
	}
	return payload, nil
}