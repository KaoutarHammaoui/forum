package routes

import (
	"net/http"

	"forum/handlers"
	"forum/middleware"
)

func Route() {
	//Home && HomeUser
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/homeUser", middleware.AuthMiddleware(handlers.HomeUser))

	//Register && login
	http.HandleFunc("/login", handlers.LoginH)
	http.HandleFunc("/do-login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogOUT)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/do-register", handlers.DoRegisterHandler)

	
}
