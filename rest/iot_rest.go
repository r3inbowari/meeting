package rest

import (
	"github.com/gorilla/mux"
	"meeting/web"
	"net/http"
)

func IotRest(r *mux.Router) {
	r.HandleFunc("/iot/list", web.PostApply).Methods(http.MethodGet) // iot列表

}
