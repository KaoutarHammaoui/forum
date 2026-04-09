package routes

import (
	"net/http"

	"forum/handlers"
//	"forum/middleware"
)

func Route() {
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/do-register", handlers.DoRegisterHandler)

	//login routes
	http.HandleFunc("/login",handlers.LoginH)
	http.HandleFunc("/do-login",handlers.LoginHandler)
	//logout
	http.HandleFunc("/logout",handlers.LogOUT)
	//middleware
	//http.HandleFunc("/",middleware.AuthMiddleware(handlers.Home))
	//also in creating posts

}
