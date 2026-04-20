package config

import (
	"forum/database"
	"log"
)

func Init() {
	database.DBinit()
	TemplateParse()
	if err := TemplateParse();err!=nil{
		log.Fatal("failed",err)
	}
}
