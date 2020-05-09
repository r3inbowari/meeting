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

func SucceedPack(w http.ResponseWriter, bangumi interface{}, total int) {
	var rq RequestResult
	rq.Data = bangumi
	rq.Total = total
	rq.Code = 0
	rq.Message = "succeed"
	jsonStr, err := json.Marshal(rq)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Fprintf(w, string(jsonStr))
}
