package mandatory_headers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IdentifiedStatusError struct {
	Status int
	Message string
}

func ByClient(c *gin.Context){
	client:=c.GetHeader("client")

	if client == "" {
		c.JSON(http.StatusBadRequest, &IdentifiedStatusError{
			Status:  http.StatusBadRequest,
			Message: "Client is required on header parameter",
		})
		c.Abort()
		return
	}
}

func GetClient(c *gin.Context) (string, *IdentifiedStatusError) {
	if c.GetHeader("client") != "" {
		return c.GetHeader("client"), nil
	}
	return "", &IdentifiedStatusError{
		Status:  http.StatusBadRequest,
		Message: "Client is required on header parameter",
	}
}