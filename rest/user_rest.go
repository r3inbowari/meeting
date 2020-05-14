package rest

import (
	"github.com/gorilla/mux"
	"meeting/web"
)

// 就是个小作业, 就不搞什么验证中间件了
func LoginRest(r *mux.Router) {
	r.HandleFunc("/login", web.Login).Methods("Post") // 登录
	r.HandleFunc("/login", web.Info).Methods("Get")   // 登录信息

	r.HandleFunc("/reg", web.Reg).Methods("Post")         // 注册
	r.HandleFunc("/reg", web.RegGet).Methods("Get")       // 注册权限
	r.HandleFunc("/reg/{uid}", web.RegPut).Methods("Put") // 注册修改

	r.HandleFunc("/user/avatar", web.Avatar).Methods("Post")   // 头像上传
	r.HandleFunc("/user/avatar", web.GetAvatar).Methods("Get") // 头像地址

	r.HandleFunc("/serviceVersion", web.ServiceVersion).Methods("Get") // 版本信息
}
