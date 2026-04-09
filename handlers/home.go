package handlers

import (
	"forum/config"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request)  {
	tmp := config.GetTemplate("home.html")
	if tmp == nil {
		//

	}
	err := tmp.Execute(w, nil)
	if err != nil {
		//
	}

}
//    background: linear-gradient(to right, rgb(5, 23, 29), rgb(0, 34, 61));
