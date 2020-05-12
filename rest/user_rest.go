package rest

import (
	"github.com/gorilla/mux"
	"meeting/web"
)

func LoginRest(r *mux.Router) {
	r.HandleFunc("/login", web.Login).Methods("Post")
	r.HandleFunc("/login", web.Info).Methods("Get")

	r.HandleFunc("/reg", web.Reg).Methods("Post")
	r.HandleFunc("/reg", web.RegGet).Methods("Get")
	r.HandleFunc("/reg/{uid}", web.RegPut).Methods("Put")

	r.HandleFunc("/user/avatar", web.Avatar).Methods("Post")
	r.HandleFunc("/user/avatar", web.GetAvatar).Methods("Get")

	r.HandleFunc("/serviceVersion", web.ServiceVersion).Methods("Get")
}
