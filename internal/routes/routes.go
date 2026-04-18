package routes

import (
	"forum/internal/handler"
	"net/http"
)

func Route() {

	
	http.HandleFunc("/", handler.Index)

}