package web

import (
	"meeting/utils"
	"net/http"
)

/**
 * 设备硬性开关
 */
func IotSwitch(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetAuthToken(r)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusUnauthorized, utils.OpAuthError)
		return
	}

	if !CheckToken(token) {
		utils.FailedResult(w, "login error", 1, http.StatusUnauthorized, utils.OpLoginError)
		return
	}
}
