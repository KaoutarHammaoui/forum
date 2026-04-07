package models

import (
	"errors"
	"forum/database"
	"time"
	"github.com/google/uuid"
)

type Session struct {
	IdSession int
	UserId    int
	Token     string
	ExpiresAt time.Time
}

//Au niveau d connexion d'un user et verification d email && pass on creer une session pour ce utilisateur afin pour chaque requete on ajoute un middlware qui verifier token qui est passe par navigateur 
func InsertSession(idUser int) (string, error) {
	query := "INSERT INTO session (user_id, token, expires_at) VALUES (?, ?, ?)"
	token := uuid.New().String()
	expires_at := time.Now().Add(time.Hour)
	_, err := database.DB.Exec(query, idUser, token, expires_at)
	if err != nil {
		return "", err
	}
	return token, nil
}

//On Récupere une Session afin d verification 
func GetSessionByToken(token string) (Session, error) {
	session := Session{}
	query := "SELECT id, user_id, token, expires_at FROM session WHERE token = ?"
	err := database.DB.QueryRow(query, token).Scan(&session.IdSession, &session.UserId, &session.Token, &session.ExpiresAt)
	if err != nil {
		return Session{}, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		DeleteSessionByToken(token)
		return Session{}, errors.New("session expirée")
	}
	return session, nil
}

//On supprime une session au niveau d logout 
func DeleteSessionByToken(token string) error {
	query := "DELETE  FROM session WHERE token = ?"
	_, err := database.DB.Exec(query, token)
	if err != nil {
		return err
	}
	return nil
}
