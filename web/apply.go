package web

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"io"
	"meeting/da"
	"meeting/utils"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

/**
 * 申请部分
 * @author r3inbowari
 * @update 2020/5/13
 */

/*
 * 教室表结构
 */
type Room struct {
	Id   string `json:"id"`   // 课室id
	Name string `json:"name"` // 课室名称
	Desc string `json:"desc"` // 介绍
}

type File struct {
	Id      string    `json:"id"`       // 文件id
	ApplyID string    `json:"apply_id"` // 申请实例id
	Name    string    `json:"name"`     // 文件名称
	Ext     string    `json:"ext"`      // 文件后缀
	Create  time.Time `json:"create"`   // 创建时间
}

type Apply struct {
	Id       string    `json:"id"`               // 申请实例id
	Uid      string    `json:"uid"`              // 用户id
	Rid      string    `json:"rid" valid:"uuid"` // 教室id
	Username string    `json:"username"`         // 申请人
	Start    time.Time `json:"start" `           // 开始时间
	End      time.Time `json:"end"`              // 结束时间
	Status   int       `json:"status"`           // 申请状态
	View     string    `json:"view"`             // 申请意见
	Content  string    `json:"content"`          // 申请内容
	Files    []File    `json:"files"`            // 文件连接
	Created  time.Time `json:"created"`          // 创建时间

}

//type ApplyRequest struct {
//	Id     string    `json:"id"`     // 课室id
//	Time   time.Time `json:"time"`   // 某一天
//	View   string    `json:"view"`   // 申请意见
//	Status int       `json:"status"` // 申请状态
//}

/**
 * 提交申请
 */
func PostApply(w http.ResponseWriter, r *http.Request) {
	// 验证
	token, err := utils.GetAuthToken(r)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusUnauthorized, utils.OpAuthError)
		return
	}

	if !CheckToken(token) {
		utils.FailedResult(w, "login error", 1, http.StatusUnauthorized, utils.OpLoginError)
		return
	}

	// 解析
	var apply Apply
	if err := utils.JsonBind(&apply, r); err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusInternalServerError, utils.OpJsonBindError)
		return
	}

	_, err = govalidator.ValidateStruct(apply)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusInternalServerError, utils.OpValidateError)
		return
	}

	// 验证时间
	if err := apply.ValidTimeCheck(); err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusOK, utils.OpValidateError)
		return
	}

	// 用户
	var dat User

	da.DBC().Where("uid = ?", GetIdTokenByToken(token)).First(&dat)

	// 用户的四个参数
	// room的id
	// 开始时间
	// 结束时间
	// 申请理由
	// 不支持跨天申请
	aid := utils.CreateUUID()
	app := Apply{
		Id:       aid,
		Uid:      dat.Uid,
		Username: dat.Username,
		Rid:      apply.Rid,
		Start:    apply.Start,
		End:      apply.End,
		Status:   0,
		Content:  apply.Content,
		Created:  time.Now(),
	}

	da.DBC().Create(app)

	utils.SucceedResult(w, "yes", 1, http.StatusOK, utils.OpSucceed)
}

/*
 * 可用性查询
 */
func (apply *Apply) ValidTimeCheck() error {
	if apply.Start.IsZero() || apply.End.IsZero() {
		return errors.New("error time")
	}

	if apply.Start.Unix() >= apply.End.Unix() {
		return errors.New("error time")
	}

	if apply.Start.Add(time.Hour*8).Day() != apply.End.Add(time.Hour*8).Day() {
		return errors.New("cross time")
	}

	for _, v := range apply.GetRoomData() {
		if apply.End.Unix() >= v.Start.Unix() && apply.End.Unix() <= v.End.Unix() {
			return errors.New("time has been occupied")
		}
		if apply.Start.Unix() >= v.Start.Unix() && apply.Start.Unix() <= v.End.Unix() {
			return errors.New("time has been occupied")
		}
		if apply.Start.Unix() >= v.Start.Unix() && apply.End.Unix() <= v.End.Unix() {
			return errors.New("time has been occupied")
		}
		if apply.Start.Unix() <= v.Start.Unix() && apply.End.Unix() >= v.End.Unix() {
			return errors.New("time has been occupied")
		}
	}

	return nil
}

/**
 * 获取某个课室的所有记录 一天或范围
 */
func (apply *Apply) GetRoomData() []Apply {
	var e []Apply
	dsb := apply.Start.Add(time.Hour * 8)
	if apply.End.Unix() == apply.Start.Unix() {
		ds := time.Date(dsb.Year(), dsb.Month(), dsb.Day(), 0, 0, 0, 0, time.Local)
		de := time.Date(dsb.Year(), dsb.Month(), dsb.Day(), 23, 59, 0, 0, time.Local)
		//var b []File
		da.DBC().Where("rid = ? AND start >= ? AND end <= ?", apply.Rid, ds, de).Find(&e)
		//da.DBC().Where("rid = ? AND start >= ? AND end <= ?", apply.Rid, ds, de).Find(&e).Related(&[]File{})
		//da.DBC().Model(&e).Related(&b)
	} else {
		dse := apply.End.Add(time.Hour * 8)
		ds := time.Date(dsb.Year(), dsb.Month(), dsb.Day(), 0, 0, 0, 0, time.Local)
		de := time.Date(dse.Year(), dse.Month(), dse.Day(), 23, 59, 0, 0, time.Local)
		da.DBC().Where("rid = ? AND start >= ? AND end <= ?", apply.Rid, ds, de).Find(&e)
	}
	return e
}

/**
 * 获取某个课室的数据
 */
func GetApply(w http.ResponseWriter, r *http.Request) {

	//if r.Method == http.MethodOptions {
	//	w.Header().Add("Access-Control-Allow-Origin", "*")
	//	w.Header().Add("content-type", "application/json")
	//	w.Header().Add("Access-Control-Allow-Headers", "Authorization")
	//	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
	//	return
	//}
	if r.Method == http.MethodOptions {
		return
	}
	// 验证
	token, err := utils.GetAuthToken(r)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusUnauthorized, utils.OpAuthError)
		return
	}

	if !CheckToken(token) {
		utils.FailedResult(w, "login error", 1, http.StatusUnauthorized, utils.OpLoginError)
		return
	}

	_ = r.ParseForm()
	var apply Apply
	if r.FormValue("rid") == "" || r.FormValue("time") == "" {
		utils.FailedResult(w, err.Error(), 1, http.StatusInternalServerError, utils.OpValidateError)
		return
	}
	apply.Rid = r.FormValue("rid")
	op, err := time.Parse("2006-01-02T15:04:05Z", r.FormValue("time"))
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusInternalServerError, utils.OpValidateError)
		return
	}
	apply.Start = op
	// 解析

	// range date data get
	if r.FormValue("end") != "" {
		ed, err := time.Parse("2006-01-02T15:04:05Z", r.FormValue("end"))
		if err != nil {
			utils.FailedResult(w, err.Error(), 1, http.StatusInternalServerError, utils.OpValidateError)
			return
		}
		apply.End = ed
	} else {
		apply.End = op
	}

	applies := apply.GetRoomData()

	for k, _ := range applies {
		if applies[k].Status == utils.RoomHasEnd {
			continue
		}
		if applies[k].End.Unix() < time.Now().Unix() {
			applies[k].Status = utils.RoomHasEnd
			da.DBC().Model(applies[k]).Update("status", utils.RoomHasEnd)
		} else if applies[k].Start.Unix() < time.Now().Unix() && applies[k].End.Unix() > time.Now().Unix() {
			applies[k].Status = utils.RoomCarryOn
		}
	}
	utils.SucceedResult(w, applies, 1, http.StatusOK, utils.OpSucceed)
}

/**
 * 获取实例文件名
 * @param apply_id 实例名
 */
func GetMeetingFileNames(w http.ResponseWriter, r *http.Request) {
	// 验证
	token, err := utils.GetAuthToken(r)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusUnauthorized, utils.OpAuthError)
		return
	}

	if !CheckToken(token) {
		utils.FailedResult(w, "login error", 1, http.StatusUnauthorized, utils.OpLoginError)
		return
	}

	var apply Apply
	var files []File
	vars := mux.Vars(r)
	da.DBC().Where("id = ?", vars["id"]).Find(&apply)
	da.DBC().Model(&apply).Related(&files)
	if GetIdTokenByToken(token) != apply.Uid || GetRoleByToken(token) != utils.RoleManager {
		utils.FailedResult(w, "access denied resources", 1, http.StatusForbidden, utils.OpResourcesDenied)
		return
	}

	utils.SucceedResult(w, files, 1, http.StatusOK, utils.OpSucceed)
}

/**
 * 申请审核部分
 */
func PutApply(w http.ResponseWriter, r *http.Request) {
	// 验证
	token, err := utils.GetAuthToken(r)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusUnauthorized, utils.OpAuthError)
		return
	}

	if !CheckToken(token) {
		utils.FailedResult(w, "login error", 1, http.StatusUnauthorized, utils.OpLoginError)
		return
	}

	if GetRoleByToken(token) != utils.RoleManager {
		utils.FailedResult(w, "access denied resources", 1, http.StatusForbidden, utils.OpResourcesDenied)
		return
	}

	err = r.ParseForm()
	if err != nil {
		utils.FailedResult(w, "wrong operation", 1, http.StatusInternalServerError, utils.OpFailed)
		return
	}
	var status int

	if r.FormValue("status") == "" {
		status = utils.RoomNotStart
	} else {
		b, err := strconv.Atoi(r.FormValue("status"))
		if err != nil {
			utils.FailedResult(w, "wrong operation", 1, http.StatusInternalServerError, utils.OpFailed)
		}
		status = b
	}

	vars := mux.Vars(r)

	da.DBC().Model(&Apply{}).Where("id = ?", vars["id"]).Updates(map[string]interface{}{"status": status, "view": r.FormValue("view")})

	utils.SucceedResult(w, "classroom apply succeed", 1, http.StatusOK, utils.OpSucceed)
}

/**
 * 文件上传
 */
func FileUpload(w http.ResponseWriter, r *http.Request) {
	// 验证
	token, err := utils.GetAuthToken(r)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusUnauthorized, utils.OpAuthError)
		return
	}

	if !CheckToken(token) {
		utils.FailedResult(w, "login error", 1, http.StatusUnauthorized, utils.OpLoginError)
		return
	}

	vars := mux.Vars(r)
	var apply Apply
	da.DBC().Where("id = ?", vars["aid"]).Find(&apply)

	if GetIdTokenByToken(token) != apply.Uid || GetRoleByToken(token) != utils.RoleManager {
		utils.FailedResult(w, "access denied to upload resources", 1, http.StatusForbidden, utils.OpResourcesDenied)
		return
	}

	uploadFile, handle, err := r.FormFile("file")
	if handle == nil {
		utils.FailedResult(w, "error handle", 1, http.StatusBadRequest, utils.OpFailed)
		return
	}

	uuid := utils.CreateUUID()
	err = os.Mkdir("./files/", 0777)
	saveFile, err := os.OpenFile("./files/"+uuid, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		utils.FailedResult(w, "file operation failed", 1, http.StatusInternalServerError, utils.OpFailed)
		return
	}
	_, _ = io.Copy(saveFile, uploadFile)

	defer uploadFile.Close()
	defer saveFile.Close()

	file := File{
		Id:      uuid,
		ApplyID: vars["aid"],
		Name:    handle.Filename,
		Create:  time.Now(),
	}
	da.DBC().Create(file)

	utils.SucceedResult(w, "upload succeed", 1, http.StatusOK, utils.OpSucceed)
}

/**
 * 教室获取
 */
func RoomList(w http.ResponseWriter, r *http.Request) {
	var rooms []Room
	da.DBC().Find(&rooms)
	utils.SucceedResult(w, rooms, len(rooms), http.StatusOK, utils.OpSucceed)
}

/**
 * 文件上传
 */
func FileDownload(w http.ResponseWriter, r *http.Request) {
	// 验证
	token, err := utils.GetAuthToken(r)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusUnauthorized, utils.OpAuthError)
		return
	}

	if !CheckToken(token) {
		utils.FailedResult(w, "login error", 1, http.StatusUnauthorized, utils.OpLoginError)
		return
	}

	vars := mux.Vars(r)
	var file File
	if da.DBC().Where("id = ?", vars["fid"]).Find(&file).RecordNotFound() {
		utils.FailedResult(w, "file not found", 1, http.StatusBadRequest, utils.OpFailed)
		return
	}

	ff, err := os.Open("./files/" + vars["fid"])
	if err != nil {
		utils.FailedResult(w, "file operation failed", 1, http.StatusInternalServerError, utils.OpFailed)
		return
	}

	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("content-disposition", "attachment; filename=\""+url.QueryEscape(file.Name)+"\"")
	// response Body
	_, err = io.Copy(w, ff)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
}
