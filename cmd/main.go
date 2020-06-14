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

/**
 * 数据库合并
 */
func CreateTable() {
	da.DBC().AutoMigrate(&web.User{})
	da.DBC().AutoMigrate(&web.Comment{})
	da.DBC().AutoMigrate(&web.Apply{})
	da.DBC().AutoMigrate(&web.File{})
	da.DBC().AutoMigrate(&web.Room{})
}
