package assync_process

import "time"

/**
DEPRECIADO
 */

type assync_operations struct {
	Id          	int64   	`gorm:"column:id;primaryKey" json:"id"`
	Hashcode 		string 		`gorm:"column:hashcode;" json:"hashcode"`
	BrokerType		string  	`gorm:"column:broker_type;" json:"broker_type"`
}

type assync_operations_items struct {
	Id            int64             `gorm:"column:id;primaryKey" json:"id"`
	Status        string            `gorm:"column:status;" json:"hashcode"`
	TimeProcess   time.Time         `gorm:"column:time_process;" json:"time_process"`
	AssyncProcess *assync_operations `gorm:"foreignkey:AssyncOperationsId;not null;"`
	AssyncOperationsId int64 		`gorm:"column:id_assync_operation;not null;" json:"id_assync_operation"`
}
