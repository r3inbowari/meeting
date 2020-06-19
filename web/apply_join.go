package web

import (
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"meeting/da"
	"meeting/utils"
	"net/http"
)

type MeetingJoin struct {
	Id   string `json:"id"`               // id
	Aid  string `json:"aid" valid:"uuid"` // 申请id
	Name string `json:"name"`             // 姓名
}

/**
 * 用户加入会议
 */
func PostJoin(w http.ResponseWriter, r *http.Request) {
	// 验证
	//token, err := utils.GetAuthToken(r)
	//if err != nil {
	//	utils.FailedResult(w, err.Error(), 1, http.StatusUnauthorized, utils.OpAuthError)
	//	return
	//}
	//
	//if !CheckToken(token) {
	//	utils.FailedResult(w, "login error", 1, http.StatusUnauthorized, utils.OpLoginError)
	//	return
	//}

	// 解析
	var join MeetingJoin
	if err := utils.JsonBind(&join, r); err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusInternalServerError, utils.OpJsonBindError)
		return
	}

	_, err := govalidator.ValidateStruct(join)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusInternalServerError, utils.OpValidateError)
		return
	}

	var e []MeetingJoin
	da.DBC().Where("aid = ?", join.Aid).Find(&e)
	if len(e) > 9 {
		utils.FailedResult(w, "person join limited", 1, http.StatusForbidden, utils.OpFailed)
		return
	}

	join.Id = utils.CreateUUID()
	da.DBC().Create(join)
	utils.SucceedResult(w, "yes", 1, http.StatusOK, utils.OpSucceed)
}

func GetJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	// 验证
	//token, err := utils.GetAuthToken(r)
	//if err != nil {
	//	utils.FailedResult(w, err.Error(), 1, http.StatusUnauthorized, utils.OpAuthError)
	//	return
	//}
	//
	//if !CheckToken(token) {
	//	utils.FailedResult(w, "login error", 1, http.StatusUnauthorized, utils.OpLoginError)
	//	return
	//}
	vars := mux.Vars(r)
	if !govalidator.IsUUIDv4(vars["id"]) {
		utils.FailedResult(w, "id is not a valid uuid", 1, http.StatusInternalServerError, utils.OpValidateError)
		return
	}

	var e []MeetingJoin
	da.DBC().Where("aid = ?", vars["id"]).Find(&e)
	utils.SucceedResult(w, e, 1, http.StatusOK, utils.OpSucceed)
}

func DelJoin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !govalidator.IsUUIDv4(vars["id"]) {
		utils.FailedResult(w, "id is not a valid uuid", 1, http.StatusInternalServerError, utils.OpValidateError)
		return
	}
	da.DBC().Where("id = ?", vars["id"]).Delete(MeetingJoin{})
	utils.SucceedResult(w, "yes", 1, http.StatusOK, utils.OpSucceed)
}



