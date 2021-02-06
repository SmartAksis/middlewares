package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/middlewares/midleware-rest/errors"
	"strings"
)

const (
	WithoutPermission = "without permission to proceed"
)

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
		restErr := errors.MethodForbidden(WithoutPermission)
		c.JSON(restErr.Status, restErr)
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
			restErr := errors.NotAuthorized("expired token")
			fmt.Print(err.Error())
			c.JSON(restErr.Status, restErr)
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
	restErr := errors.MethodForbidden(WithoutPermission)
	result := GetClaim(claims, "scope")

	if result != nil {
		for _, v := range result.([]interface{}) {
			if v.(string) == scope {
				return
			}
		}
		c.JSON(restErr.Status, restErr)
		c.Abort()
		return
	}
}

func checkAuthorities(c *gin.Context, claims jwt.MapClaims, roles ... string) {
	restErr := errors.MethodForbidden("without permission to proceed")
	result := GetClaim(claims, "authorities")

	if roles != nil && result == nil {
		c.JSON(restErr.Status, restErr)
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
		c.JSON(restErr.Status, restErr)
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
