package model

type Session struct {
	Id     uint   `gorm:"primaryKey"`
	UserId int64  `gorm:"index"`
	User   User   `gorm:"foreignKey:UserId"`
	Token  string `gorm:"index;varchar(300)"`
}
