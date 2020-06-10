package rest

import (
	"github.com/gorilla/mux"
	"meeting/web"
)

func IotRest(r *mux.Router) {
	r.HandleFunc("/apply", web.PostApply).Methods("Post") // 申请教室
}
