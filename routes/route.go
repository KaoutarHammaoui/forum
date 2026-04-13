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
	http.HandleFunc("/add-comment", middleware.AuthMiddleware(handlers.AddCommentHandler))
	http.HandleFunc("/react", middleware.AuthMiddleware(handlers.ReactionHandler))
	//also in creating posts

}
