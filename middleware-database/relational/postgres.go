package relational

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	postgresDb *gorm.DB
)

func GetPostgresDatabase(property GormProperties) *gorm.DB {
	if postgresDb == nil {
		_DSN:=fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d search_path=%s sslmode=%s TimeZone=%s", property.Host, property.User, property.Password, property.DbName, property.Port, property.Schema, property.SslMode, property.TimeZone)
		dialect := postgres.New(postgres.Config{
			DSN: _DSN,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		})

		db, err := gorm.Open(dialect, &gorm.Config{})
		if err != nil {
			panic("Error to connect with database")
			return nil
		}
		postgresDb = db
	}


	return postgresDb
}