package relational

import (
	"gorm.io/gorm"
	"os"
)

var ParameterDao ParameterDaoInterface = &parameterDao{
	db: GetPostgresDatabase(GormDefaultConfig(os.Getenv("SCHEMA_DEFAULT"))),
}

type ParameterDaoInterface interface {
	GetParameters(key string) ([]string, error)
}

type parameterDao struct {
	db *gorm.DB
}

func (this parameterDao) GetParameters(key string) ([]string, error) {
	var result []string
	applicationKey := &ParameterApplicationKey{}
	err := this.db.Model(&ParameterApplicationKey{}).Where("key = ? AND enable = ?", key, true).Find(&applicationKey).Error
	if err != nil {
		return nil, err
	}
	if applicationKey != nil && applicationKey.Id == 0 {
		return nil, nil
	}

	var applicationValues []ParameterApplicationValue
	errValues := this.db.Model(&ParameterApplicationValue{}).Where("id_key = ? AND enable = ?", applicationKey.Id, true).Find(&applicationValues).Error

	if errValues != nil {
		return nil, errValues
	}

	if applicationValues != nil {
		result = make([]string, len(applicationValues))
		for index, value := range applicationValues {
			result[index]=value.Val
		}

	}

	return result, err
}

type ParameterApplicationKey struct {
	Id          	int64   	`gorm:"column:id;primaryKey" 	json:"id"`
	Key       		string 		`gorm:"column:key;" 			json:"key"`
	Enable 			bool 		`gorm:"column:enable;" 			json:"enable"`
}


type ParameterApplicationValue struct {
	Id          	int64   	`gorm:"column:id;primaryKey" 	json:"id"`
	IdKey       	string 		`gorm:"column:id_key;" 			json:"id_key"`
	Val		       	string 		`gorm:"column:val;" 			json:"val"`
	Enable 			bool 		`gorm:"column:enable;" 			json:"enable"`
}

