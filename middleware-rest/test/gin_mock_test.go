package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)



type ginMock struct {
	_body interface{}
	_queryString map[string]string
	_url string
}

func GinMockBuilder() ginMock {
	return ginMock{}
}

func (e ginMock) Body(body interface{}) ginMock {
	e._body = body
	return e
}

func (e ginMock) QueryStringMap(queryMap map[string]string) ginMock {
	e._queryString = queryMap
	return e
}

func (e ginMock) Url(url string) ginMock {
	e._url = url
	return e
}

func (e ginMock) Build() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	for name, value := range e._queryString {
		c.Params = []gin.Param{gin.Param{Key: name, Value: value}}
	}
	if e._body != nil {
		jsonValue, _ := json.Marshal(e._body)
		c.Request, _ = http.NewRequest(http.MethodPost, e._url, bytes.NewBuffer(jsonValue))
	}
	return c, w
}
