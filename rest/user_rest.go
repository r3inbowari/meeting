package rest

import (
	"github.com/gorilla/mux"
	"meeting/web"
)

func LoginRest(r *mux.Router) {
	r.HandleFunc("/login", web.Login).Methods("Post")
	r.HandleFunc("/reg", web.Reg).Methods("Post")
	r.HandleFunc("/reg", web.RegGet).Methods("Get")
	r.HandleFunc("/reg/{uid}", web.RegPut).Methods("Put")

	r.HandleFunc("/user/avatar", web.Avatar).Methods("Put")
}
