package routes

import (
	"net/http"

	"forum/handlers"
	"forum/middleware"
)

func Route() {
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/do-register", handlers.DoRegisterHandler)

	//login routes
	http.HandleFunc("/login",handlers.LoginH)
	http.HandleFunc("/do-login",handlers.LoginHandler)
	//logout
	http.HandleFunc("/logout",handlers.LogOUT)

	// http.HandleFunc("/",handlers.HomeHAndler)
	//middleware
	http.HandleFunc("/",middleware.AuthMiddleware(handlers.HomeHAndler))
	//also in creating posts

}
