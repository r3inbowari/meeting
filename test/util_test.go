package test

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"meeting/utils"
	"testing"
)

type Foo struct {
	A string `valid:"ipv4"`
	B string `valid:"mac"`
	C int    `valid:"range(0|100)"`
}

type bar struct {
	X   string `valid:"ipv4"`
	Foo `valid:",required"`
}

func TestJsonValid(t *testing.T) {
	utils.InitValidator()
	b := bar{
		X: "192.168.1.1",
	}

	b.Foo.A = "192.168.1.1.1"
	b.Foo.B = "aa:bb:cc:dd:ee:ff"
	b.Foo.C = 100

	result, err := govalidator.ValidateStruct(b)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println(result)
}
