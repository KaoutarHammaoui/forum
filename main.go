package main

import (
	"forum/database"
	"log"
	"net/http"
)

func main(){
	database.DBinit()
	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080",nil))

}