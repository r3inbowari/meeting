package web

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"meeting/da"
	"meeting/utils"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	Uid      string `json:"uid" gorm:"primary_key" valid:"alphanum,length(5|12)" `
	Username string `json:"username" valid:"required,length(3|16)"`
	Password string `json:"password,omitempty" valid:"ascii,length(3|16)"`
	Sex      int    `json:"sex"`
	Role     string `json:"role"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"status"`
}

type UserInfo struct {
	Uid      string `json:"uid" valid:"alphanum,length(5|12)"`
	Password string `json:"password" valid:"ascii,length(3|16)"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var ui UserInfo
	if err := utils.JsonBind(&ui, r); err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusInternalServerError, utils.OpJsonBindError)
		return
	}

	_, err := govalidator.ValidateStruct(ui)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusInternalServerError, utils.OpValidateError)
		return
	}

	token, err := LoginUser(ui)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusUnauthorized, utils.OpLoginError)
		return
	}
	utils.SucceedResult(w, token, 1, http.StatusOK, utils.OpSucceed)
}

func Reg(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := utils.JsonBind(&u, r); err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusInternalServerError, utils.OpJsonBindError)
		return
	}

	_, err := govalidator.ValidateStruct(u)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusInternalServerError, utils.OpValidateError)
		return
	}

	err = LogonUser(u)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusOK, utils.OpLogonError)
		return
	}

	utils.SucceedResult(w, "wait for audit", 1, http.StatusOK, utils.OpSucceed)
}

func RegGet(w http.ResponseWriter, r *http.Request) {
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

	us := GetAuditUserSet()
	utils.SucceedResult(w, us, len(us), http.StatusOK, utils.OpSucceed)

}

func RegPut(w http.ResponseWriter, r *http.Request) {
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

	var role string
	var status int
	if r.FormValue("role") == "" {
		role = utils.RoleUser
	} else {
		role = r.FormValue("role")
	}
	if r.FormValue("status") == "" {
		status = utils.StatusNormal
	} else {
		b, err := strconv.Atoi(r.FormValue("status"))
		if err != nil {
			utils.FailedResult(w, "wrong operation", 1, http.StatusInternalServerError, utils.OpFailed)
		}
		status = b
	}

	vars := mux.Vars(r)
	da.DBC().Model(&User{}).Where("uid = ?", vars["uid"]).Updates(map[string]interface{}{"role": role, "status": status})

	utils.SucceedResult(w, "account audit succeed", 1, http.StatusOK, utils.OpSucceed)
}

/**
 * 设置头像
 */
func Avatar(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetAuthToken(r)
	if err != nil {
		utils.FailedResult(w, err.Error(), 1, http.StatusUnauthorized, utils.OpAuthError)
		return
	}

	if !CheckToken(token) {
		utils.FailedResult(w, "login error", 1, http.StatusUnauthorized, utils.OpLoginError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		utils.FailedResult(w, "wrong operation", 1, http.StatusInternalServerError, utils.OpFailed)
		return
	}

	url := r.FormValue("url")
	da.DBC().Model(&User{}).Where("uid = ?", GetIdTokenByToken(token)).Updates(map[string]interface{}{"avatar": url})

	utils.SucceedResult(w, "account avatar change succeed", 1, http.StatusOK, utils.OpSucceed)
}

func (u *User) GetUsername() string {
	return u.Username
}

/**
 * 获取所有等待审核的账号
 */
func GetAuditUserSet() []User {
	var users []User
	da.DBC().Where("status = ?", utils.StatusAudit).Find(&users)
	for k, _ := range users {
		users[k].Password = ""
	}
	return users
}

/**
 * 审核申请和权限
 * @param uid 被审核人
 * @param role 权限
 */
func LogonConfirm(uid, role string) {
	var data User
	if !da.DBC().Where("uid = ?", uid).First(&data).RecordNotFound() {
		data.Status = utils.StatusNormal
		if role == utils.RoleUser || role == utils.RoleManager {
			data.Role = role
		} else {
			data.Role = utils.RoleUser
		}
		da.DBC().Save(&data)
	}
}

/**
 * 注册用户请求
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
			Password: utils.CreateMD5(u.Password),
			Sex:      u.Sex,
			Role:     utils.RoleDrift,
			Avatar:   u.Avatar,
			Status:   utils.StatusAudit,
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
		// 状态判断
		if userDb.Status == utils.StatusLock {
			return "", errors.New("your account has been locked by system")
		} else if userDb.Status == utils.StatusAudit {
			return "", errors.New("your account is waiting for audit")
		}
		if userDb.Password == utils.CreateMD5(user.Password) {
			return CreateToken(userDb.Uid), nil
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
 * 用户锁定检查 弃用
 */
func GetUserStatus(token string) int {
	result, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return utils.GetConfig().GetJwtSecret(), nil
	})
	if err != nil {
		if err.Error() == "Token is expired" && result != nil {
			utils.Info("token is expired", logrus.Fields{"uid": result.Claims.(jwt.MapClaims)["jti"]})
		} else {
			utils.Warn("parse with claims failed")
		}
		return utils.StatusUndefined
	}

	var data User
	if !da.DBC().Where("uid = ?", result.Claims.(jwt.MapClaims)["jti"]).First(&data).RecordNotFound() {
		// 此处role是必要的
		return data.Status
	} else {
		return utils.StatusUndefined
	}
}

/**
 * 用户鉴权
 */
func GetRoleByToken(token string) string {
	result, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return utils.GetConfig().GetJwtSecret(), nil
	})
	if err != nil {
		if err.Error() == "Token is expired" && result != nil {
			utils.Info("token is expired", logrus.Fields{"uid": result.Claims.(jwt.MapClaims)["jti"]})
		} else {
			utils.Warn("parse with claims failed")
		}
		return utils.RoleDrift
	}

	var data User
	if !da.DBC().Where("uid = ?", result.Claims.(jwt.MapClaims)["jti"]).First(&data).RecordNotFound() {
		// 此处role是必要的
		return data.Role
	} else {
		return utils.RoleDrift
	}
}

/**
 * 获取用户id
 */
func GetIdTokenByToken(token string) string {
	result, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return utils.GetConfig().GetJwtSecret(), nil
	})
	if err != nil {
		if err.Error() == "Token is expired" && result != nil {
			utils.Info("token is expired", logrus.Fields{"uid": result.Claims.(jwt.MapClaims)["jti"]})
		} else {
			utils.Warn("parse with claims failed")
		}
		return ""
	}

	var data User
	if !da.DBC().Where("uid = ?", result.Claims.(jwt.MapClaims)["jti"]).First(&data).RecordNotFound() {
		// 此处role是必要的
		return data.Uid
	} else {
		return ""
	}

}

/**
 * 创建token
 * @param uid 用户id
 */
func CreateToken(uid string) string {
	claims := &jwt.StandardClaims{
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
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
