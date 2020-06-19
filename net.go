package meeting

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"meeting/rest"
	"meeting/utils"
	"net/http"
	"strconv"
	"time"
)

/**
 * 启动位置
 * @author r3inbowari
 * @create 5/28/2020
 *
 * @history 6/11/2020 add iot rest uri
 *          5/29/2020 created default options
 *          6/20/2020 sql runtime
 */

func Start() {
	r := mux.NewRouter()
	rest.LoginRest(r)
	rest.CommentRest(r)
	rest.ApplyRest(r)
	rest.IotRest(r)
	r.Use(corsMid) // 中间件拦截
	utils.Info(utils.GetConfig().Name + " service port: 9999")

	err := http.ListenAndServe(":9999", r)
	if err != nil {
		log.Println("[INFO] 监听端口失败, 五秒后退出")
		time.Sleep(time.Second * 5)
		return
	}
}

func corsMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()                                                //获取本地现在时间
		w.Header().Add("Access-Control-Allow-Origin", "*")              //允许访问所有域
		w.Header().Add("content-type", "application/json")              //返回数据格式是json
		w.Header().Add("Access-Control-Allow-Headers", "Authorization") //header的类型
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		w.Header().Add("Access-Control-Max-Age", "600")
		next.ServeHTTP(w, r)
		t2 := time.Now()
		d := t2.Sub(t1) //两个时间相减
		utils.Info("mux route", logrus.Fields{"path": r.URL.Path, "method": r.Method, "runtime": strconv.FormatInt(d.Milliseconds(), 10) + "ms"})
	})
}
