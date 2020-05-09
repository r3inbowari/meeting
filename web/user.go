package web

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"meeting/utils"
	"net/http"
	"time"
)

type User struct {
	Uid      string `json:"uid"`
	Username string `json:"username"`
	Password string `json:"password"`
	Sex      int    `json:"sex"`
	GroupID  string `json:"groupID"`
	Avatar   string `json:"avatar"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	utils.SucceedPack(w, "hello", 1)
}

func (u *User) GetUsername() string {
	return u.Username
}

type UserInfo struct {
	Uid      string `json:"uid"`
	Password string `json:"password"`
}

func CheckLogin(ui UserInfo) bool {
	return false
}

/**
 * 检查token
 * @param token
 */
func CheckToken(token string) bool {
	result, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return utils.GetConfig().GetJwtSecret(), nil
	})
	if err != nil {
		if err.Error() == "Token is expired" && result != nil {
			utils.Info("token is expired", logrus.Fields{"uid": result.Claims.(jwt.MapClaims)["jti"]})
		} else {
			utils.Warn("parse with claims failed")
		}
		return false
	}
	return true
}

/**
 * 创建token
 * @param uid 用户id
 */
func CreateToken(uid string) string {
	claims := &jwt.StandardClaims{
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * 5).Unix(),
		Issuer:    "r3inb",
		Id:        uid,
	}
	utils.Info("create token "+uid, logrus.Fields{"uid": uid})
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(utils.GetConfig().GetJwtSecret())
	if err != nil {
		println(err)
		return ""
	}
	return ss
}
