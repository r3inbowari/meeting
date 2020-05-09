package test

import (
	"meeting"
	"meeting/da"
	"meeting/web"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	da.DBC().AutoMigrate(&web.User{})
	time.Sleep(time.Second * 1)

	da.DBC().Create(web.User{
		Uid:      "admin",
		Username: "admin",
		Password: "admin",
		Sex:      meeting.Famale,
		GroupID:  "ab1",
		Avatar:   "http://icon.qq.com",
	})
	time.Sleep(time.Second * 5)
}
