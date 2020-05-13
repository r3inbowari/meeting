package rest

import (
	"github.com/gorilla/mux"
	"meeting/web"
)

func FeedbackRest(r *mux.Router) {
	r.HandleFunc("/feedback", web.Feedback).Methods("Post")
}
