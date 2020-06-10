package test

import (
	"meeting/da"
	"meeting/utils"
	"meeting/web"
	"testing"
	"time"
)

func TestApply(t *testing.T) {

	da.DBC().AutoMigrate(&web.Apply{})
	da.DBC().AutoMigrate(&web.File{})
	time.Sleep(time.Second)
	da.DBC().Model(&web.File{}).AddForeignKey("apply_id", "apply(id)", "RESTRICT", "RESTRICT")
	aid := utils.CreateUUID()

	app := web.Apply{
		Id:      aid,
		Uid:     "admin",
		Rid:     utils.CreateUUID(),
		Start:   time.Now(),
		End:     time.Now().Add(time.Hour),
		Status:  0,
		Content: "chenhaibin",
		Created: time.Now(),
		Files: []web.File{
			{
				Id:      utils.CreateUUID(),
				ApplyID: aid,
				Name:    "kuaikuai.doc",
				Create:  time.Now(),
			},
			{
				Id:      utils.CreateUUID(),
				ApplyID: aid,
				Name:    "xiaobimssd.jpg",
				Create:  time.Now(),
			},
		},
	}

	da.DBC().Create(app)
}
