package routes

import (
	"forum/internal/handler"
	"forum/internal/middleware"
	"net/http"
)

func Route() {
	http.HandleFunc("/", middleware.CheckUserContext(handler.Index))
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/logout", handler.LogOUT)
	http.HandleFunc("/posts", middleware.CheckUserContext(middleware.AuthMiddleware(handler.CreatePost)))
	http.HandleFunc("/comments", middleware.CheckUserContext(middleware.AuthMiddleware(handler.CreateComment)))
	http.HandleFunc("/reactions", middleware.CheckUserContext(middleware.AuthMiddleware(handler.ReactPost)))
	http.HandleFunc("/static/", handler.StaticHandler)
	http.HandleFunc("/uploads/", handler.UploadsHandler)
}
