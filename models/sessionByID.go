package models
 import(
	"forum/database"
 )
func DeleteSessionsByUserID(userID int) error {
	_, err := database.DB.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	return err
}
