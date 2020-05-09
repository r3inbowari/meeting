package meeting

import (
	"log"
	"meeting/rest"
	"net/http"
	"time"
)

func Start()  {
	rest.LoginRest()
	log.Println("[INFO] 解析服务启动: 9999")
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Println("[INFO] 监听端口失败, 五秒后退出")
		time.Sleep(time.Second * 5)
		return
	}
}
