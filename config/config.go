package config

import "time"

const (
	//Server Config
	APP_HOST string = "" //default "" is localhost
	APP_PORT string = "8000"

	//JWT Config
	JWT_EXP        time.Duration = 24 //hours
	JWT_SECRET_KEY string        = "21fd35418dbcbb27e69857dcf015c20a7179feb9d89d2efed05fe465815bf9cc"

	//Database Connection Config
	DB_HOST string = "127.0.0.1"
	DB_PORT string = "5432"
	DB_USER string = "postgres"
	DB_PASS string = "secret123"
	DB_NAME string = "db_gopos"
	DB_PROV string = "postgres" // available 2 provider : mysql & postgres

	//Database Connection Pool
	DB_POOL_IDLE int           = 10
	DB_POOL_OPEN int           = 100
	DB_POOL_TIME time.Duration = 1 * time.Hour
)
