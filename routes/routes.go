package routes

import (
	"net/http"

	"forum/handlers"
	"forum/middleware"
)

func Route() {

	//Assets Routes :
	http.HandleFunc("/static/", handlers.StaticHandler)
	http.HandleFunc("/uploads/", handlers.UploadsHandler)

	//Home Route :
	http.HandleFunc("/", handlers.Home)

	//Register Routes:
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/do-register", handlers.DoRegisterHandler)


	// Login Routes:
	http.Handle("/login", middleware.RateLimiter(http.HandlerFunc(handlers.LoginH)))
	http.HandleFunc("/do-login", handlers.LoginHandler)
	// Logout Route:
	http.HandleFunc("/logout", handlers.LogOUT)

	//Home User && Create Posts Routes : 
	http.Handle("/homeUser",
		middleware.RateLimiter(middleware.AuthMiddleware(http.HandlerFunc(handlers.HomeUser))))

	http.Handle("/createPost",
		middleware.RateLimiter(middleware.AuthMiddleware(http.HandlerFunc(handlers.CreatePost))))

	//comments
	http.Handle("/SubmitComment",middleware.RateLimiter(middleware.AuthMiddleware(http.HandlerFunc(handlers.AddComment))))	
}
