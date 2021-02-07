package database

import (
	"fmt"
	"github.com/middlewares/middleware-tracing/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	database *gorm.DB
)

func GetDatabase(property GormProperties) *gorm.DB {
	if database == nil {

		dialect := postgres.New(postgres.Config{
			DSN: fmt.Sprintf("host=%d user=%d password=%d dbname=%d port=%d search_path=%d sslmode=%d TimeZone=%d", property.Host, property.User, property.Password, property.DbName, property.Port, property.Schema, property.SslMode, property.TimeZone),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		})

		db, err := gorm.Open(dialect, &gorm.Config{})
		if err != nil {
			logger.Error("Error to get connection", err)
			return nil
		}

		database = db

	}


	return database
}