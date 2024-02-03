package database

import (
	"log"

	"github.com/GoPOS-id/gopos-api/api/model"
	"github.com/GoPOS-id/gopos-api/config"
	"gorm.io/gorm"
)

func Init() {

	database := dbConfig(config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASS, config.DB_NAME, config.DB_PROV)
	db, err := database.connect()

	if err != nil {
		log.Fatalln(err)
	}
	// db.Logger = db.Logger.LogMode(logger.Silent)

	DB = db
	database.connectionPool(config.DB_POOL_IDLE, config.DB_POOL_OPEN, config.DB_POOL_TIME)

	db.AutoMigrate(
		&model.User{},
		&model.Session{},
		&model.Role{},
	)
}

func DbContext() *gorm.DB {
	return DB
}
