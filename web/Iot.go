package web

import (
	"meeting/utils"
	"net/http"
	"time"
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
	Id     string    `json:"id"`               // 物联网uuid
	Name   string    `json:"name"`             // 名称
	Rid    string    `json:"rid" valid:"uuid"` // 关联教室id
	Type   int       `json:"type"`             // 设备类型
	Create time.Time `json:"create"`           // 创建时间
}

/**
 * 设备列表获取
 * @param rid 指向的room_id
 * @history 11/6/2020 create
 */
func IotList(w http.ResponseWriter, r *http.Request) {
	names := GetDTUSessionsKey()
	utils.SucceedResult(w, names, 1, http.StatusOK, utils.OpSucceed)
}

/**
 *IotOrder
 */
//func IotList(w http.ResponseWriter, r *http.Request) {
//	names := GetDTUSessionsKey()
//	utils.SucceedResult(w, names, 1, http.StatusOK, utils.OpSucceed)
//}