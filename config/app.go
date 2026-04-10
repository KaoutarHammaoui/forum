package config

import (
	"forum/database"
	"log"
)

func Init() {
	// Configuration de l'application, comme la connexion à la base de données, les variables d'environnement, etc.
	database.DBinit()
	LoadAssets()
	TemplateParse()
	if err := TemplateParse();err!=nil{
		log.Fatal("failed",err)
	}
}
