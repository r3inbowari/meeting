package main

import (
	"meeting"
	"meeting/utils"
	"time"
)

func main() {
	utils.InitLogger()
	meeting.Start()


	time.Sleep(time.Second * 50)



}

