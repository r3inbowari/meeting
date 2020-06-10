package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestResult struct {
	Total   int         `json:"total"`
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"msg"`
}

func SucceedResult(w http.ResponseWriter, data interface{}, total int, tag int, code int) {
	var rq RequestResult
	rq.Data = data
	rq.Total = total
	rq.Code = code
	rq.Message = "succeed"
	jsonStr, err := json.Marshal(rq)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	w.WriteHeader(tag)

	println(string(jsonStr))
	fmt.Fprintf(w, string(jsonStr))
}

func FailedResult(w http.ResponseWriter, data interface{}, total int, tag int, code int) {
	var rq RequestResult
	rq.Data = data
	rq.Total = total
	rq.Code = code
	rq.Message = "failed"
	jsonStr, err := json.Marshal(rq)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	println(string(jsonStr))
	w.WriteHeader(tag)

	fmt.Fprintf(w, string(jsonStr))
}
