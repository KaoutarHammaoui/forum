package routes

import (
	"net/http"

	"forum/handlers"
)

func Route() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/do-register", handlers.DoRegisterHandler)
}
