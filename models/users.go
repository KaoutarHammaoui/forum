package models

import (
	"time"

	"forum/database"
)

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

// Insertion d'un NV user au niveau d registrer :(Kuhaku)
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

// Verifier Email et Username si ils sont dupliquée au niveau d registrer un nv user (kuuhaku)
func ExistsInColumn(column, value string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users  WHERE " + column + " = ?"
	err := database.DB.QueryRow(query, value).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Au niveau d login on est besoin de check User exists et d'apres ca tu valide email, password, username !!!!! (Kaoutar)
func GetUserByEmail(email string) (User, error) {
	payload := User{}
	query := "SELECT id, username, email, password FROM users  WHERE email = ?"
	err := database.DB.QueryRow(query, email).Scan(&payload.ID, &payload.Username, &payload.Email, &payload.Password)
	if err != nil {
		return User{}, err
	}
	return payload, nil
}

// on est besoin de cette func au niveau  :
// Verification d utilisateur est il connecté a une relation au niveau d session
// Affichage d Posts ....
func GetUserByID(id int) (User, error) {
	payload := User{}
	query := "SELECT id, username, email, password, created_at  FROM users  WHERE id = ?"
	err := database.DB.QueryRow(query, id).Scan(&payload.ID, &payload.Username, &payload.Email, &payload.Password, &payload.CreatedAt)
	if err != nil {
		return User{}, err
	}
	return payload, nil
}
