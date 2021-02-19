package relational

import (
	"os"
	"strconv"
)

type GormProperties struct {
	Host string
	User string
	Password string
	Port int
	DbName string
	TimeZone string
	SslMode string
	Schema string
}

func GormDefaultConfig(schema string) GormProperties {
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	return GormProperties{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     port,
		DbName:   os.Getenv("DB_NAME"),
		TimeZone: os.Getenv("DB_TIME_ZONE"),
		SslMode:  "disable",
		Schema:   schema,
	}
}
