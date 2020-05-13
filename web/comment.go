package web

import (
	"meeting/da"
	"meeting/utils"
	"net/http"
	"time"
)

type Comment struct {
	Content string    `json:"content" gorm:"type:varchar(100)"`
	Uid     string    `json:"uid"`
	Type    int       `json:"type"`
	Attach  string    `json:"attach"`
	Created time.Time `json:"created"`
}

func Feedback(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetAuthToken(r)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusUnauthorized, utils.OpAuthError)
		return
	}

	if !CheckToken(token) {
		utils.FailedResult(w, "login error", 1, http.StatusUnauthorized, utils.OpLoginError)
		return
	}

	var dat User
	da.DBC().Where("uid = ?", GetIdTokenByToken(token)).First(&dat)

	var feed Comment
	if err := utils.JsonBind(&feed, r); err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusInternalServerError, utils.OpJsonBindError)
		return
	}

	feed.Type = utils.ContentFeedback
	feed.Attach = "feedback"
	feed.Uid = dat.Uid
	feed.Created = time.Now()

	da.DBC().Create(&feed)
	utils.SucceedResult(w, "pull content succeed", 1, http.StatusOK, utils.OpSucceed)
}
