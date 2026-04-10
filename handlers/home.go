package handlers

import (
	"forum/config"
	"net/http"
)

func HomeHAndler(w http.ResponseWriter, r *http.Request) {
	tmpl:=config.GetTemplate("home.html")
	if tmpl==nil{
		http.Error(w,"not found",http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w,"home.html",nil)
}
