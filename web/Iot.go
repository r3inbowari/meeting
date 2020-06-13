package web

import (
	"net/http"
)

/**
 * 物联网设备
 * @author r3inbowari
 * @create 6/11/2020
 *
 */

/**
 * 硬性开关
 * @method get
 */
func IotSwitch(w http.ResponseWriter, r *http.Request) {
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

}

type Iot struct {
	Id  string `json:"id"`               // 申请实例id
	Rid string `json:"rid" valid:"uuid"` // 教室id

}

/**
 * 设备列表获取
 * @param rid 指向的room_id
 * @history 11/6/2020 create
 */
func IotList(w http.ResponseWriter, r *http.Request) {

}

/**
 *
 */
