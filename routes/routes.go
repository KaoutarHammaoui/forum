package routes

import (
	"net/http"

	"forum/handlers"
	"forum/middleware"
)

func Route() {

	// Assets Routes
	http.HandleFunc("/static/", handlers.StaticHandler)

	// Home Route
	http.HandleFunc("/", handlers.Home)

	// Register Routes
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.Handle("/do-register", middleware.RateLimiter(http.HandlerFunc(handlers.DoRegisterHandler)))

	// Login Routes
	http.HandleFunc("/login", handlers.LoginH)
	http.Handle("/do-login", middleware.RateLimiter(http.HandlerFunc(handlers.LoginHandler)))

	// Logout Route
	http.HandleFunc("/logout", handlers.LogOUT)

	// Home User && Create Posts Routes
	http.Handle("/homeUser", middleware.RateLimiter(middleware.AuthMiddleware(http.HandlerFunc(handlers.HomeUser))))
	http.Handle("/createPost", middleware.RateLimiter(middleware.AuthMiddleware(http.HandlerFunc(handlers.CreatePost))))

	// Comments
	http.Handle("/reactions", middleware.RateLimiter(middleware.AuthMiddleware(handlers.ReactPost)))
	http.Handle("/SubmitComment", middleware.RateLimiter(middleware.AuthMiddleware(http.HandlerFunc(handlers.AddComment))))
	http.Handle("/reactOnAComment", middleware.RateLimiter(middleware.AuthMiddleware(http.HandlerFunc(handlers.ReactComment))))
}
