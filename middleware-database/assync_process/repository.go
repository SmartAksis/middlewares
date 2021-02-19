package assync_process

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"os"
	"time"
)

/**
DEPRECIADO
 */

const (
	SqsBrokerAssyncType = "SQS"
	SentStatusType = "SENT"
)

var (
	database *gorm.DB
)

func nextValOperation(sequenceName string) int64{
	var assyncProcessId int64
	stmt:=fmt.Sprintf("SELECT nextval('%s.\"%s\"')", os.Getenv("SCHEMA_DEFAULT"), sequenceName)
	database.Raw(stmt).Scan(&assyncProcessId)
	return assyncProcessId;
}

func Consumed(messageId string, queueName string){
	initDatabase()
}
func InitSqsAssyncProcess(messageId string, queueName string) (error) {
	initDatabase()
	operationId:=nextValOperation("seq_assync_operations")
	err := database.Save(&assync_operations{
		Id:         operationId,
		Hashcode:   messageId,
		BrokerType: SqsBrokerAssyncType,
	}).Error
	if err != nil{
		return err
	} else {
		statusErr := database.Save(&assync_operations_items{
			Id:            nextValOperation("seq_assync_opt_status"),
			Status:        SentStatusType,
			TimeProcess:   time.Now(),
			AssyncOperationsId: operationId,
		}).Error

		if statusErr != nil {
			return statusErr
		}

	}
	return nil
}

func initDatabase() {
	if database == nil {
		Host:= os.Getenv("DB_HOST")
		User:= os.Getenv("DB_USER")
		Password:= os.Getenv("DB_PASSWORD")
		Port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
		DbName:= os.Getenv("DB_NAME")
		TimeZone:= os.Getenv("DB_TIME_ZONE")
		SslMode:= "disable"
		Schema:= os.Getenv("SCHEMA_DEFAULT")


		dialect := postgres.New(postgres.Config{
			DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d search_path=%s sslmode=%s TimeZone=%s", Host, User, Password, DbName, Port, Schema, SslMode, TimeZone),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		})

		db, err := gorm.Open(dialect, &gorm.Config{})
		if err == nil {
			database = db
		}
	}
}