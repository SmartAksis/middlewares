package relational

import (
	"fmt"
	"github.com/smart-aksis/golang-middlewares/middleware-rest/request_utils"
)

func GetFilter(filters ...request_utils.FilterField) (string, []interface{}) {
	var parameters []interface{}
	query := ""
	if len(filters) > 0 {
		for _, filter := range filters {
			if len(parameters) == 0 {
				query+=fmt.Sprintf("%s = ? ", filter.Field)
			} else {
				query+=fmt.Sprintf("%s %s = ? ", filter.Operation, filter.Field)
			}
			parameters = append(parameters, filter.Value)
		}
	}
	return query, parameters
}

func Paginate() {
	//err = this.db.Model(&entities.User{}).Where(relational.GetFilter(filters...)).Find(&users).Error
}
