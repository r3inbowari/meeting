package main

import (
	"meeting"
	"meeting/da"
	"meeting/utils"
	"meeting/web"
	"time"
)

func main() {
	CreateTable()
	utils.InitLogger()
	meeting.Start()

	time.Sleep(time.Second * 50)
}

func CreateTable() {
	da.DBC().AutoMigrate(&web.User{})
	da.DBC().AutoMigrate(&web.Comment{})
}
