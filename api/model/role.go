package model

type Role struct {
	Id   uint   `gorm:"primaryKey"`
	Name string `gorm:"index"`
	User []User
}
