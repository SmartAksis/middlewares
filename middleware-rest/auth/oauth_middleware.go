package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	WithoutPermission = "without permission to proceed"
)

type ResponseErrorToken struct {
	Message string
	Error	string
}

func Authenticated(c *gin.Context){
	token, _ := GetToken(c)
	if token == nil {
		c.Abort()
	}
}

func AuthenticatedRead(c *gin.Context) {
	token, claims := GetToken(c)
	if token == nil {
		c.Abort()
		return
	}
	if claims == nil {
		c.Abort()
		return
	}
	authenticateScope(c, claims, "read")
}


func AuthenticatedWrite(c *gin.Context) {
	token, claims := GetToken(c)
	if token == nil {
		c.Abort()
		return
	}
	if claims == nil {
		c.Abort()
		return
	}
	authenticateScope(c, claims, "write")
}

func authenticateScope(c *gin.Context, claims jwt.MapClaims, _scope string){
	result := GetClaim(claims, "scope")
	if result != nil {
		for _, v := range result.([]interface{}) {
			if v.(string) == _scope {
				return
			}
		}
		c.JSON(http.StatusForbidden, &ResponseErrorToken{
			Message: WithoutPermission,
			Error: "forbidden",
		})
		c.Abort()
		return
	}
}


func GetToken(c *gin.Context) (*jwt.Token, jwt.MapClaims) {
	request := c.Request
	if request == nil {
		return nil, nil
	}
	authorization := request.Header.Get("Authorization")
	if authorization != "" && strings.HasPrefix(authorization, "Bearer ") {
		tokenValue := strings.Replace(authorization, "Bearer ", "", 1)
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("smartaksiskey"), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, &ResponseErrorToken{
				Message: "expired token",
				Error: "unauthorized",
			})
			return nil, nil
		}
		return token, claims
	} else {
		return nil, nil
	}
	return nil, nil
}

func CheckPermissions(c *gin.Context, scope string, roles ... string){
	token, claims := GetToken(c)
	if token != nil && claims != nil {
		checkAuthorities(c, claims, roles...)
		checkScopes(c, claims, scope)
	} else {
		c.Abort()
	}
}

func CheckAuthorities(c *gin.Context, roles ... string) {
	token, claims := GetToken(c)
	if token != nil && claims != nil {
		checkAuthorities(c, claims, roles...)
	} else {
		c.Abort()
	}
}

func CheckScopes(c *gin.Context, scope string) {
	token, claims := GetToken(c)
	if token != nil && claims != nil {
		checkScopes(c, claims, scope)
	} else {
		c.Abort()
	}}

func checkScopes(c *gin.Context, claims jwt.MapClaims, scope string) {
	result := GetClaim(claims, "scope")

	if result != nil {
		for _, v := range result.([]interface{}) {
			if v.(string) == scope {
				return
			}
		}
		c.JSON(http.StatusForbidden, &ResponseErrorToken{
			Message: WithoutPermission,
			Error: "forbidden",
		})
		c.Abort()
		return
	}
}

func checkAuthorities(c *gin.Context, claims jwt.MapClaims, roles ... string) {
	result := GetClaim(claims, "authorities")

	if roles != nil && result == nil {
		c.JSON(http.StatusForbidden, &ResponseErrorToken{
			Message: "without permission to proceed",
			Error: "forbidden",
		})
		c.Abort()
		return
	}

	if result != nil {
		for _, v := range result.([]interface{}) {
			_, boolean := Find(roles, v.(string))
			if boolean {
				return
			}
		}
		c.JSON(http.StatusForbidden, &ResponseErrorToken{
			Message: "without permission to proceed",
			Error: "forbidden",
		})
		c.Abort()
		return
	}

}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func GetClaim(claims jwt.MapClaims, _key string) interface{} {
	for key, val := range claims {
		if key == _key {
			return val
		}
	}
	return nil
}
