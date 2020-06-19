package rest

import (
	"github.com/gorilla/mux"
	"meeting/web"
	"net/http"
)

func IotRest(r *mux.Router) {
	r.HandleFunc("/iot/list", web.IotList).Methods(http.MethodGet) // iot列表

	//r.HandleFunc("/iot/list", web.IotOrder).Methods(http.MethodGet) // iot列表
}
