package test

import (
	"meeting"
	"meeting/da"
	"meeting/web"
	"testing"
	"time"
)

func TestAddUser(t *testing.T) {
	// create
	da.DBC().AutoMigrate(&web.User{})
	time.Sleep(time.Second * 1)
	data := web.User{
		Uid:      "admin",
		Username: "admin",
		Password: "admin",
		Sex:      meeting.Famale,
		GroupID:  "ab1",
		Avatar:   "http://icon.qq.com",
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
		Sex:      meeting.Famale,
		GroupID:  "ab1",
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
	}

	time.Sleep(time.Second * 4)

}
