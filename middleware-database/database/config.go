package database

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
	return GormProperties{
		Host:     "localhost",
		User:     "smart_aksis",
		Password: "sm4r74k515",
		Port:     5432,
		DbName:   "postgres",
		TimeZone: "Brazil/East",
		SslMode:  "disable",
		Schema:   schema,
	}
}
