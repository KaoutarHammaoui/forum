package routes

import (
	"net/http"

	"forum/handlers"
	"forum/middleware"
)

func Route() {
	http.HandleFunc("/", handlers.Home)
	// login routes
	http.Handle("/login", middleware.RateLimiter(http.HandlerFunc(handlers.LoginH)))
	http.HandleFunc("/do-login", handlers.LoginHandler)
	// logout
	http.HandleFunc("/logout", handlers.LogOUT)

	http.Handle("/homeUser", 
		middleware.RateLimiter(middleware.AuthMiddleware(http.HandlerFunc(handlers.HomeUser),),),)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/do-register", handlers.DoRegisterHandler)
}
