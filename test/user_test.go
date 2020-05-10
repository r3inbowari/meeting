package test

import (
	"encoding/json"
	"meeting/da"
	"meeting/utils"
	"meeting/web"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	a := web.CreateToken("restore")
	println(a)
	println(web.CheckToken(a))

	time.Sleep(time.Second * 6)
	println(web.CheckToken(a))
}

func TestAddUser(t *testing.T) {
	// create
	da.DBC().AutoMigrate(&web.User{})
	time.Sleep(time.Second * 1)
	data := web.User{
		Uid:      "admin",
		Username: "admin",
		Password: utils.CreateMD5("admin"),
		Sex:      utils.Famale,
		Role:     "manager",
		Avatar:   "http://icon.qq.com",
		Status:   utils.StatusNormal,
	}

	// insert
	da.DBC().Create(&data)
	time.Sleep(time.Second * 10)
}

func TestOperationUser(t *testing.T) {
	// create
	da.DBC().AutoMigrate(&web.User{})
	time.Sleep(time.Second * 1)
	data := web.User{
		Uid:      "admin",
		Username: "admin",
		Password: "admin",
		Sex:      utils.Famale,
		Role:     "ab1",
		Avatar:   "http://icon.qq.com",
	}

	// insert
	da.DBC().Create(&data)
	time.Sleep(time.Second * 10)

	var dat web.User
	// select one warn: danger operation
	da.DBC().Where("uid = ?", "admin").First(&dat)
	println(dat.Username)
	time.Sleep(time.Second * 5)

	// del
	da.DBC().Delete(&dat)
	time.Sleep(time.Second * 5)

}

func TestUserExist(t *testing.T) {
	_, exist := web.FoundUser("admin")
	println(exist)
}

func TestLogonAndLogin(t *testing.T) {
	_ = web.LogonUser(web.User{Uid: "r3inb", Password: "15598870762", Username: "chen"})

	// 模拟审核过程
	web.LogonConfirm("r3inb", "manager")
	time.Sleep(time.Second * 2)

	ui := web.UserInfo{
		Uid:      "r3inb",
		Password: "15598870762",
	}
	token, err := web.LoginUser(ui)
	if err != nil {
		println(err.Error())
	} else {
		println(token)
		a := web.GetRoleByToken(token)
		println(a)
	}
	time.Sleep(time.Second * 4)
}

func TestGetLockUserSet(t *testing.T) {
	a := web.GetAuditUserSet()
	println(len(a))
	b, _ := json.Marshal(a)
	println(string(b))
}
