package database

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var DB*sql.DB

func DBinit(){
	var err error
	DB,err=sql.Open("sqlite3","./forum.db")
	if err !=nil{
		log.Fatal(err)
	}
	if err = DB.Ping();err!=nil{
        log.Fatal(err)
		
	}
	DB.Exec("PRAGMA foreign_keys = ON")

	TableCreation()

}