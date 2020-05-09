package rest

import (
	"meeting/web"
	"net/http"
)

func LoginRest() {
	http.HandleFunc("/login", web.Login)

}


