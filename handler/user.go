package handler

import (
	"filestore-server/db"
	"filestore-server/util"
	"net/http"
)

const (
	pwdSalt = string("#*2021")
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Redirect(w, r, "./static/view/signup", http.StatusFound)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(username) < 3 || len(password) < 5 {
		w.WriteHeader(http.StatusBadRequest)
	}

	encPassword := util.Sha1([]byte(password+pwdSalt))

	suc := db.UserSignup(username, encPassword)

	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}

}

