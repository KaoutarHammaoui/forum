package main

import (
	"log"
	"net/http"

	"forum/config"
	"forum/routes"
)

func main() {
	config.Init()
	routes.Route()

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
