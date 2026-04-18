package routes

import (
	"forum/internal/handler"
	"forum/internal/middleware"
	"net/http"
)

func Route() {
	http.HandleFunc("/", middleware.CheckUserContext(handler.Index))
}