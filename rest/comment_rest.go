package rest

import (
	"github.com/gorilla/mux"
	"meeting/web"
)

func CommentRest(r *mux.Router) {
	r.HandleFunc("/feedback", web.Feedback).Methods("Post")
}
