package routes

import (
	"forum/internal/handler"
	"forum/internal/middleware"
	"net/http"
)

func Route() {
	http.HandleFunc("/", middleware.CheckUserContext(handler.Index))
	http.HandleFunc("/login", handler.Login)
	// http.HandleFunc("/register", handler.register)

}