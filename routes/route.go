package routes

import (
	"forum/handlers"
	"forum/middleware"
	"net/http"
)

func Route() {
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/do-register", handlers.DoRegisterHandler)
	http.Handle("/comment", middleware.AuthMiddleware(http.HandlerFunc(handlers.AddCommentHandler)))
	http.Handle("/react", middleware.AuthMiddleware(http.HandlerFunc(handlers.ReactionHandler)))
}
