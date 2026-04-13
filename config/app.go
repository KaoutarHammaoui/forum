package config

import (
	"forum/database"
	"log"
)

func Init() {
	database.DBinit()
	LoadAssets()
	// FIX: was called twice — once silently (result ignored) then again with error check
	if err := TemplateParse(); err != nil {
		log.Fatal("failed to parse templates: ", err)
	}
}
