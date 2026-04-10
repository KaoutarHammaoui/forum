package handlers

 import (
 	"forum/config"
 	"forum/models"
 	"net/http"
 )

 type Data struct {
 	Posts []models.Post
 	Categories []models.Category

 }
 func HomeUser(w http.ResponseWriter, r *http.Request) {

 	 //sécurité méthode
 	if r.Method != http.MethodGet && r.Method != http.MethodPost {
 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
 		return
 	}

 	 //POST = création de post
 	if r.Method == http.MethodPost {

 		err := r.ParseMultipartForm(10 << 20) // 10 MB max
 		if err != nil {
 			http.Error(w, "Invalid form data", http.StatusBadRequest)
 			return
 		}

 		title := r.FormValue("title")
 		content := r.FormValue("content")
 		categories := r.Form["categories"]

 		// image upload (optionnel)
 		file, handler, err := r.FormFile("image")
 		var imagePath string

 		if err == nil {
 			defer file.Close()

 			 //ici tu peux stocker dans /uploads
 			imagePath = "uploads/" + handler.Filename
 			// TODO: save file physically
 		}

 		//post := models.Post{}
 		 //création post
 		//err = models.InsertPost(post)
 		if err != nil {
 			http.Error(w, "Error creating post", http.StatusInternalServerError)
 			return
 		}

 		// IMPORTANT : rester sur la même page
 		http.Redirect(w, r, "/homeUser", http.StatusSeeOther)
 		return
 	}

 	 //GET = affichage page
 	data := Data{}

 	categories, err := models.GetAllCategory()
 	if err != nil {
 		http.Error(w, "Error loading categories", http.StatusInternalServerError)
 		return
 	}
 	data.Categories = categories

 	posts, err := models.GetAllPosts()
 	if err != nil {
 		http.Error(w, "Error loading posts", http.StatusInternalServerError)
 		return
 	}
 	data.Posts = posts

 	tmpl := config.GetTemplate("homeUser.html")
 	if tmpl == nil {
 		http.Error(w, "Template not found", http.StatusInternalServerError)
 		return
 	}

 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
 	_ = tmpl.Execute(w, data)
 }