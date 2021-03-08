package auth

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

const (
	WithoutPermission = "without permission to proceed"
)

type ResponseErrorToken struct {
	Message string
	Error	string
}

func BasicReadAuthenticated(c *gin.Context) {
	readClient:=os.Getenv("SMART_AKSIS_READ_CLIENT")
	readPass:=os.Getenv("SMART_AKSIS_READ_PASS")
	if readClient == "" || readPass == "" {
		c.JSON(forbidden("There's no configuration to basic read authentication"))
		c.Abort()
		return
	}
	sEnc := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", readClient, readPass)))

	request := c.Request
	if request == nil {
		c.JSON(forbidden("There's no request"))
		c.Abort()
		return
	}
	authorization := request.Header.Get("Authorization")

	if authorization != "" && strings.HasPrefix(authorization, "Basic ") {
		authorizationValue := strings.Replace(authorization, "Basic ", "", 1)
		if sEnc != authorizationValue {
			c.JSON(forbidden("There's no authentication"))
			c.Abort()
			return
		}
	}
}

func getResponseErrorToken(httpStatus int, message string, error string) (int, ResponseErrorToken){
	return httpStatus, ResponseErrorToken{
		Message: message,
		Error: error,
	}
}

func forbidden(message string) (int, ResponseErrorToken){
	return getResponseErrorToken(http.StatusUnauthorized, message, "forbidden")
}

func unauthorized(message string)(int, ResponseErrorToken){
	return getResponseErrorToken(http.StatusForbidden, message, "unauthorized")
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
		c.JSON(forbidden("Without jwt"))
		c.Abort()
		return
	}
	if claims == nil {
		c.JSON(forbidden("Without claims"))
		c.Abort()
		return
	}
	authenticateScope(c, claims, "read")
}


func AuthenticatedWrite(c *gin.Context) {
	token, claims := GetToken(c)
	if token == nil {
		c.JSON(forbidden("Without jwt"))
		c.Abort()
		return
	}
	if claims == nil {
		c.JSON(forbidden("Without claims"))
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
		c.JSON(forbidden(WithoutPermission))
		c.Abort()
		return
	}
}


func GetToken(c *gin.Context) (*jwt.Token, jwt.MapClaims) {
	signKey := os.Getenv("JWT_SIGN_KEY")
	request := c.Request
	if request == nil {
		return nil, nil
	}
	authorization := request.Header.Get("Authorization")
	if authorization != "" && strings.HasPrefix(authorization, "Bearer ") {
		tokenValue := strings.Replace(authorization, "Bearer ", "", 1)
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(signKey), nil
		})
		if err != nil {
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
		c.JSON(forbidden(WithoutPermission))
		c.Abort()
		return
	}
}

func checkAuthorities(c *gin.Context, claims jwt.MapClaims, roles ... string) {
	result := GetClaim(claims, "authorities")

	if roles != nil && result == nil {
		c.JSON(forbidden("without permission to proceed"))
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
		c.JSON(forbidden("without permission to proceed"))
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
