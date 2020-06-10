package rest

import (
	"github.com/gorilla/mux"
	"meeting/web"
	"net/http"
)

func ApplyRest(r *mux.Router) {
	r.HandleFunc("/apply", web.PostApply).Methods("Post") // 申请教室
	r.HandleFunc("/apply", web.GetApply).Methods(http.MethodGet, http.MethodOptions)   // 申请教室

	r.HandleFunc("/meeting/file/{aid}", web.FileUpload).Methods("Post")  // 会议文件上传
	r.HandleFunc("/meeting/file/{fid}", web.FileDownload).Methods("Get") // 会议文件下载

	r.HandleFunc("/meeting/filenames/{id}", web.GetMeetingFileNames).Methods("Get") // 教室文件名获取

	r.HandleFunc("/apply/{id}", web.PutApply).Methods("Put") // 审核修改

}
