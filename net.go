package meeting

import (
	"github.com/gorilla/mux"
	"log"
	"meeting/rest"
	"net/http"
	"time"
)

func Start()  {
	r := mux.NewRouter()
	rest.LoginRest(r)
	log.Println("[INFO] 解析服务启动: 9999")
	err := http.ListenAndServe(":9999", r)
	if err != nil {
		log.Println("[INFO] 监听端口失败, 五秒后退出")
		time.Sleep(time.Second * 5)
		return
	}
}
