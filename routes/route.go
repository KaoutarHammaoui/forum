package routes

import (
	"net/http"

	"forum/handlers"
)

func Route() {
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/do-register", handlers.DoRegisterHandler)

	//login routes
	http.HandleFunc("/login",handlers.LoginH)
	http.HandleFunc("/do-login",handlers.LoginHandler)
	//logout
	http.HandleFunc("/logout",handlers.LogOUT)
}
