package request_utils

import (
	"encoding/json"
	"fmt"
)

	func MapFromObject(data interface{}) (map[string]interface{}, error) {
	converted, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var mapResult map[string]interface{}
	err = json.Unmarshal(converted, &mapResult)
	return mapResult, err
}

func QueryFromObject(data interface{}) (string, error) {
	if data == nil {
		return "", nil
	}
	mapResult, err:=MapFromObject(data)
	if err != nil {
		return "", err
	}
	query:=""
	for key, value := range mapResult {
		if value != ""{
			if query == "" {
				query+="?"
			} else {
				query+="&"
			}
			var resultValue string
			resultValue = ConvertData(value)
			if key != "id" {
				query+=fmt.Sprintf("%s=%s", key, resultValue)
			} else {
				if resultValue != "0" {
					query+=fmt.Sprintf("%s=%s", key, resultValue)
				}
			}
		}
	}
	return query, nil
}
