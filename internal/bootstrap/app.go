package bootstrap

import (
	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/routes"
	"log"
)

func Init() {
	database.DataBaseinit()
	if err := config.TemplateParse(); err != nil {
		log.Fatal(err)
	}
	routes.Route()
}
