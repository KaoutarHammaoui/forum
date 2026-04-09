package handlers

import (
	"net/http"
	"strings"
	"time"

	"forum/config"
	"forum/models"

	"golang.org/x/crypto/bcrypt"
)

func LoginH(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tmpl := config.GetTemplate("login.html")
	if tmpl == nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	tmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return

	}
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	data := Login{Email: email}

	user, err := models.GetUserByEmail(email)
	if err != nil {
		data.EmailError = "no accoubnt"
		data.HasErrors = true
		renderLogin(w,data)
		return
	}

	if err:=bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password));err!=nil{
		data.PassError="error i pass"
		data.HasErrors=true
		renderLogin(w,data)
		
		return
	}

	token,err:=models.InsertSession(user.ID)
	if err!=nil{
		http.Error(w,"",http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,&http.Cookie{
		Name: "token",
		Value: token,
		Expires: time.Now().Add(1* time.Hour),
		HttpOnly: true,
		Path: "/",
	})
	http.Redirect(w,r,"/",http.StatusSeeOther)
}

func LogOUT(w http.ResponseWriter, r *http.Request){
	
	cookie,err:=r.Cookie("token")

	if err==nil{
		models.DeleteSessionByToken(cookie.Value)
	}
	http.SetCookie(w,&http.Cookie{
		Name: "token",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path: "/",
	})
	http.Redirect(w,r,"/login",http.StatusSeeOther)
}

func renderLogin(w http.ResponseWriter, data Login) {
	tmpl := config.GetTemplate("login.html")
	if tmpl == nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
	
}
