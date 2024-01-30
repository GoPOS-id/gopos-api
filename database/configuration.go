package database

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type databaseContext struct {
	Hostname string
	Port     string
	Username string
	Password string
	Database string
	Provider string
}

func (db *databaseContext) dsn() string {
	if db.Provider == "mysql" {
		return db.mysqlDsn()
	}
	return db.postgresDsn()
}

func (db *databaseContext) mysqlDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", db.Username, db.Password, db.Hostname, db.Port, db.Database)
}

func (db *databaseContext) postgresDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", db.Hostname, db.Username, db.Password, db.Database, db.Port)
}

func (d *databaseContext) connect() (*gorm.DB, error) {
	var database *gorm.DB
	if d.Provider == "mysql" {
		var db, err = gorm.Open(mysql.Open(d.dsn()), &gorm.Config{
			SkipDefaultTransaction: true,
		})

		if err != nil {
			return &gorm.DB{}, err
		}
		database = db
	} else {
		var db, err = gorm.Open(postgres.Open(d.dsn()), &gorm.Config{
			SkipDefaultTransaction: true,
		})

		if err != nil {
			return &gorm.DB{}, err
		}
		database = db
	}

	return database, nil
}

// connectionPool sets up and returns the connection pool configurations.
func (d *databaseContext) connectionPool(idle int, open int, time time.Duration) *sql.DB {
	conn, _ := DB.DB()
	conn.SetMaxIdleConns(idle)
	conn.SetMaxOpenConns(open)
	conn.SetConnMaxLifetime(time)
	return conn
}

// dbConfig creates and returns a databaseType configuration.
func dbConfig(hostname string, port string, username string, password string, database string, provider string) databaseContext {
	return databaseContext{
		Hostname: hostname,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
		Provider: provider,
	}
}
