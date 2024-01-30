package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id         int64  `gorm:"primaryKey"`
	Username   string `gorm:"type:varchar(25)"`
	Password   string `gorm:"varchar(300)"`
	Fullname   string `gorm:"type:varchar(30)"`
	Email      string `gorm:"index"`
	RoleId     uint
	Role       Role       `gorm:"foreignKey:RoleId"`
	VerifiedAt *time.Time `gorm:"index"`
	CreatedAt  time.Time
	UpdateAt   time.Time
	DeletedAt  gorm.DeletedAt
	Session    []Session
}
