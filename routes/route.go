package routes

import (
	"net/http"

	"forum/handlers"
	"forum/middleware"
)

func Route() {
	http.HandleFunc("/", handlers.Home)
	//login routes
	http.HandleFunc("/login", handlers.LoginH)
	http.HandleFunc("/do-login", handlers.LoginHandler)
	//logout
	http.HandleFunc("/logout", handlers.LogOUT)

	http.HandleFunc("/homeUser", middleware.AuthMiddleware(handlers.HomeUser))
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/do-register", handlers.DoRegisterHandler)

}
