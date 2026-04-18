package cmd

import (
	"fmt"
	"forum/internal/bootstrap"
	"log"
	"net/http"
)


func Serve(){
	bootstrap.Init()
	log.Println("http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}