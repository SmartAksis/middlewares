package data

import (
	"encoding/json"
	"github.com/smart-aksis/golang-middlewares/middleware-rest/request_utils"
	"gorm.io/gorm"
)

type GenericDaoInterface interface {
	GetModel() (tx *gorm.DB)
}

func Paginate(dao GenericDaoInterface, filters []request_utils.FilterField, paginationProperties request_utils.PaginationProperties, dest interface{}) error {
	var result *gorm.DB
	if len(filters) > 0 {
		for _, filter := range filters {
			if result == nil {
				result = dao.GetModel().Where(filter.Field, filter.Value)
			} else {
				if filter.Operation == "OR" {
					result.Or(filter.Field, filter.Value)
				}
				if filter.Operation != "OR" {
					result.Where(filter.Field, filter.Value)
				}
			}
		}
	} else {
		result = dao.GetModel()
	}

	var limit int
	var page int

	limit = paginationProperties.PageSize
	page = paginationProperties.PageNumber - 1

	if page < 0 {
		page = 0
	}

	if limit < 1 {
		limit = 10
	}

	return result.Limit(limit).Offset(page * limit).Find(dest).Error
}

func ConvertToMapByJsonDefinitions (data interface{}) (map[string]interface{}, error) {
	converted, err := json.Marshal(data)
	if err != nil {
		return nil, err;
	}
	var mapResult map[string]interface{}
	err = json.Unmarshal(converted, &mapResult)
	return mapResult, err
}
