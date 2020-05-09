package web

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"meeting/da"
	"meeting/utils"
	"net/http"
	"time"
)

type User struct {
	Uid      string `json:"uid" gorm:"primary_key" valid:"required"`
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
	Sex      int    `json:"sex"`
	GroupID  string `json:"groupID"`
	Avatar   string `json:"avatar"`
}

type UserInfo struct {
	Uid      string `json:"uid" valid:"required"`
	Password string `json:"password" valid:"required"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	utils.SucceedPack(w, "hello", 1)
}

func (u *User) GetUsername() string {
	return u.Username
}

/**
 * 注册用户
 * @param uid 用户id
 * @param password 密码
 * @return error login info
 */
func LogonUser(u User) error {
	_, err := govalidator.ValidateStruct(u)
	if err != nil {
		return err
	}

	if _, exist := FoundUser(u.Uid); exist {
		user := User{
			Uid:      u.Uid,
			Username: u.Username,
			Password: u.Password,
			Sex:      u.Sex,
			GroupID:  u.GroupID,
			Avatar:   u.Avatar,
		}
		da.DBC().Create(&user)
		return nil
	} else {
		return errors.New("the uid has been occupied")
	}
}

/**
 * 登录用户
 * @param user UserInfo
 * @return string token
 * @param error login info
 */
func LoginUser(user UserInfo) (string, error) {
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		return "", err
	}

	if userDb, exist := FoundUser(user.Uid); exist {
		return "", errors.New("non existing user")
	} else {
		if userDb.Password == utils.CreateMD5(user.Password) {
			return CreateToken(user.Uid), nil
		}
		return "", errors.New("user password error")
	}
}

/**
 * 用户是否存在
 * @param uid 用户id
 */
func FoundUser(uid string) (User, bool) {
	var user User
	if da.DBC().Where("uid = ?", uid).First(&user).RecordNotFound() {
		return user, true
	}
	return user, false
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
