package main

import (
	"github.com/sirupsen/logrus"
	"meeting"
	"meeting/da"
	"meeting/utils"
	"meeting/web"
	"net"
	"strconv"
	"time"
)

func main() {
	CreateTable()
	utils.InitLogger()
	go RunDTUService()
	meeting.Start()
	time.Sleep(time.Second * 50)
	// esp8266 tcp透传
}

func TCPService() {

}

const Network = "tcp"

var listener net.Listener

func RunDTUService() {
	var err error
	port := utils.GetConfig().IotPort
	if port == nil {
		port = new(int)
		*port = 8000
	}
	if listener, err = net.Listen(Network, ":"+strconv.Itoa(*port)); err != nil {
		utils.Fatal("iot listen failed", logrus.Fields{"port": *port, "err": err})
	} else {
		utils.Info("iot listened ok", logrus.Fields{"port": *port})
	}
	defer func() { _ = listener.Close() }()
	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.Warn("accept failed", logrus.Fields{"err": err})
			break
		}
		go web.DtuHandle(conn)
	}
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
