package utils

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlConnPool(connectionString string) *gorm.DB {

	// create connection string to mysql
	// todo make class handle many database
	db, err := gorm.Open(mysql.Open(connectionString))

	// just panic error
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	// config connection pool
	mysqlDB, _ := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	mysqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	mysqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	mysqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
