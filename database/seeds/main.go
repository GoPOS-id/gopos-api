package main

import (
	"fmt"
	"time"

	"github.com/GoPOS-id/gopos-api/app/model"
	"github.com/GoPOS-id/gopos-api/constant"
	"github.com/GoPOS-id/gopos-api/database"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	database.Init()
	db := database.DbContext()

	Role := []model.Role{
		{Id: 1, Name: constant.Role_Operator},
		{Id: 2, Name: constant.Role_Adminstrator},
		{Id: 3, Name: constant.Role_Cashier},
	}

	if err := db.Create(&Role).Error; err != nil {
		fmt.Println("Error seeds Role: " + err.Error())
		return
	}

	passwd, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	currentTime := time.Now()
	pointerTime := &currentTime
	guid := xid.New()
	User := model.User{
		Id:         guid.Time().Unix(),
		Username:   "admin123",
		Password:   string(passwd),
		Fullname:   "Operator Account",
		Email:      "operator@gopos.com",
		RoleId:     1,
		VerifiedAt: pointerTime,
	}

	if err := db.Create(&User).Error; err != nil {
		fmt.Println("Error seeds User: " + err.Error())
		return
	}

	fmt.Println("Success seeds starter account! ❤️")
}
