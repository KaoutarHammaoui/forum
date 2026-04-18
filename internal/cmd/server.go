package cmd

import (
	"fmt"
	"log"
	"net/http"
)


func Serve(){
	log.Println("http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("")
	}
}