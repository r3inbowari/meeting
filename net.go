package meeting

import (
	"github.com/gorilla/mux"
	"log"
	"meeting/rest"
	"net/http"
	"time"
)

func Start() {
	r := mux.NewRouter()
	rest.LoginRest(r)
	rest.CommentRest(r)
	rest.ApplyRest(r)
	rest.IotRest(r)
	r.Use(corsMid) // 跨域处理
	log.Println("[INFO] 解析服务启动: 9999")
	err := http.ListenAndServe(":9999", r)
	if err != nil {
		log.Println("[INFO] 监听端口失败, 五秒后退出")
		time.Sleep(time.Second * 5)
		return
	}
}

func corsMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		w.Header().Add("Access-Control-Allow-Origin", "*")              //允许访问所有域
		w.Header().Add("content-type", "application/json")              //返回数据格式是json
		w.Header().Add("Access-Control-Allow-Headers", "Authorization") //header的类型
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")

		next.ServeHTTP(w, r)
	})
}
