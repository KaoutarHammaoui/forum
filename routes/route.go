package routes

import (
	"net/http"

	"forum/handlers"
)

func Route() {
	http.HandleFunc("/", handlers.Home)
	//login routes
	http.HandleFunc("/login", handlers.LoginH)
	http.HandleFunc("/do-login", handlers.LoginHandler)
	//logout
	http.HandleFunc("/logout", handlers.LogOUT)


	
	http.HandleFunc("/homeUser", handlers.HomeUser)
	//http.HandleFunc("/", handlers.HomeUser)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/do-register", handlers.DoRegisterHandler)

	



}
