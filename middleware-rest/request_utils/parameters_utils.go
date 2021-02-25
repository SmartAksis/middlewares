package request_utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RequestError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

type PaginationProperties struct {
	PageNumber int `json:"page"`
	PageSize int `json:"size"`
}

type FilterField struct {
	Field 		string 		`json:"pageNumber"`
	Operation 	string 		`json:"operation"`
	Value 		interface{} `json:"value"`
}


func ConvertData(value interface{}) string {

	chain:=&stringConverter{
		next: &float64Converter{
			next: &float32Converter{
				next: &float32Converter{
					next: &int64Converter{
						next: &int32Converter{
							next:&int16Converter{
								next:&int8Converter{
									next:&intConverter{

									},
								},
							},
						},
					},
				},
			},
		},
	}
	return chain.Convert(value)
}

func convertValue(value string) interface{}{
	intValue, err := strconv.Atoi(value)
	if err == nil {
		return intValue
	}
	float64Value, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return float64Value
	}
	float32Value, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return float32Value
	}
	return value
}

func FilterFieldAnd(field string, value string) FilterField {
	return FilterField{ Field: field, Operation: "AND", Value: convertValue(value)}
}

func FilterFieldOr(field string, value string) FilterField {
	return FilterField{ Field: field, Operation: "OR", Value: convertValue(value) }
}

func FilterFieldLike(field string, value string) FilterField {
	return FilterField{ Field: field, Operation: "LIKE", Value: convertValue(value) }
}


func badRequestError(message string) *RequestError {
	return &RequestError{
		Message: message,
		Status: http.StatusBadRequest,
		Error: "Invalid request",
	}
}

func PathNumberInVariable(c *gin.Context, key string) (int64, *RequestError) {
	numberParameter:=c.Params.ByName(key)
	if numberParameter == "" {
		return 0, badRequestError("Id parameter is required")
	}

	number, err:=strconv.Atoi(numberParameter)
	if err != nil {
		return 0, badRequestError("Id parameter is invalid")
	}

	if number < 0 {
		return 0, badRequestError("Id parameter is negative")
	}

	return int64(number), nil
}

func SimpleQueryFilter(c *gin.Context, keys ...string) []FilterField  {
	var result []FilterField
	if keys != nil && len(keys) > 0 {
		for _, value := range keys {
			if c.Query(value) != "" {
				result=append(result, FilterFieldAnd(value, c.Query(value)))
			}
		}
	}
	return result
}

func GetPaginationParameter(c *gin.Context) PaginationProperties{
	c.Query("page")
	var pojo PaginationProperties
	if err := c.ShouldBindJSON(&pojo); err != nil {
		return PaginationProperties{
			PageNumber: 0,
			PageSize:   10,
		}
	}
	if pojo.PageSize == 0 {
		pojo.PageSize = 10
	}
	return pojo
}

func GetPageNumber(c *gin.Context) int{
	c.Query("size")
	var pojo PaginationProperties
	if err := c.ShouldBindJSON(&pojo); err != nil {
		return 0
	}
	return pojo.PageNumber
}

func GetPageSize(c *gin.Context) int{
	c.Query("page")
	var pojo PaginationProperties
	if err := c.ShouldBindJSON(&pojo); err != nil {
		return 10
	}
	return pojo.PageNumber
}

