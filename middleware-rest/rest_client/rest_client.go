package rest_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type IntegrationStatus struct {
	Status 	int
	Method 	string
	Url 	string
	Message string
	Error	error
}

type RestConsumer interface {
	BaseEndpoint() string
}

type endPointEnv struct {
	_host string
	_port string
	_ssl string
	_baseUrl string
}

func EndPointEnvBuilder() endPointEnv {
	return endPointEnv{}
}

func (e endPointEnv) Build() string {
	var protocol string
	var host string
	var port int
	var baseUrl string

	if e._ssl != "" && os.Getenv(e._ssl) != "" {
		protocol = os.Getenv(e._ssl)
	} else {
		protocol = "http"
	}

	if e._host != "" {
		host = os.Getenv(e._host)
	}

	if e._port == "" || os.Getenv(e._port) == ""{
		panic("Error to localizate port to client")
	}
	port, _ = strconv.Atoi(os.Getenv(e._port))

	if e._baseUrl != "" {
		baseUrl = os.Getenv(e._baseUrl)
	}

	if baseUrl != "" {
		return fmt.Sprintf("%s://%s:%d/%s", protocol, host, port, baseUrl)
	} else {
		return fmt.Sprintf("%s://%s:%d", protocol, host, port)
	}

}

func (e endPointEnv) Host(host string) endPointEnv {
	e._host = host
	return e
}

func (e endPointEnv) Port(port string) endPointEnv {
	e._port = port
	return e
}

func (e endPointEnv) Ssl(ssl string) endPointEnv {
	e._ssl = ssl
	return e
}

func (e endPointEnv) BaseUrl(baseUrl string) endPointEnv {
	e._baseUrl = baseUrl
	return e
}


func Get(rest RestConsumer, url string, typeVar interface{}, headerMap map[string]string) IntegrationStatus {
	client := &http.Client{}
	uri:=rest.BaseEndpoint() + url
	request, err := http.NewRequest(http.MethodGet, uri, nil)
	for index, value := range headerMap {
		request.Header.Set(index, value)
	}

	resp, err := client.Do(request)
	method:="GET"

	if err != nil {
		return IntegrationStatus{ Message: resp.Status, Status: resp.StatusCode, Method: method, Url: uri, Error: err,}
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return IntegrationStatus{ Message: resp.Status, Status: resp.StatusCode, Method: method, Url: uri, Error: err,}
	}

	if resp.StatusCode >= 400 {
		return IntegrationStatus{ Message: resp.Status, Status: resp.StatusCode, Method: method, Url: uri, Error: errors.New("Error in ResultApi"),}
	}

	jsonUnMarshallError:=json.Unmarshal(body, typeVar)
	if jsonUnMarshallError != nil {
		return IntegrationStatus{ Message: resp.Status, Status: resp.StatusCode, Method: method, Url: uri, Error: jsonUnMarshallError,}
	}

	return IntegrationStatus{
		Status:  resp.StatusCode,
		Method:  method,
		Url:     uri,
		Message: resp.Status,
		Error:   nil,
	}
}


func Post(rest RestConsumer, url string, body map[string]string, typeVar interface{}, headerMap map[string]string) IntegrationStatus {
	jsonValue, _ := json.Marshal(body)
	client := &http.Client{}
	uri:=rest.BaseEndpoint() + url

	var request *http.Request
	var err error

	if body != nil {
		request, err = http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonValue))
	} else {
		request, err = http.NewRequest(http.MethodPost, uri, nil)
	}

	for index, value := range headerMap {
		request.Header.Set(index, value)
	}

	resp, err := client.Do(request)
	method:="POST"

	if err != nil {
		return IntegrationStatus{ Status: resp.StatusCode, Method: method, Url: uri, Error: err,}
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return IntegrationStatus{ Status: resp.StatusCode, Method: method, Url: uri, Error: err,}
	}

	if resp.StatusCode >= 400 {
		return IntegrationStatus{ Status: resp.StatusCode, Method: method, Url: uri, Error: err,}
	}

	jsonUnMarshallError:=json.Unmarshal(responseBody, typeVar)
	if jsonUnMarshallError != nil {
		return IntegrationStatus{ Status: resp.StatusCode, Method: method, Url: uri, Error: jsonUnMarshallError,}
	}

	return IntegrationStatus{
		Status:  resp.StatusCode,
		Method:  method,
		Url:     uri,
		Message: resp.Status,
		Error:   nil,
	}
}



func Patch(rest RestConsumer, url string, body map[string]interface{}, typeVar interface{}) IntegrationStatus {
	client := &http.Client{}
	jsonValue, _ := json.Marshal(body)
	uri:=rest.BaseEndpoint() + url
	//stringConv:=fmt.Sprintf("%v", body)
	request, err := http.NewRequest(http.MethodPatch, uri, bytes.NewBuffer(jsonValue))
	response, err := client.Do(request)

	if err != nil {
		return IntegrationStatus{ Status: response.StatusCode, Method: http.MethodPatch, Url: uri, Error: err,}
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return IntegrationStatus{ Status: response.StatusCode, Method: http.MethodPatch, Url: uri, Error: err,}
	}

	if response.StatusCode >= 400 {
		return IntegrationStatus{ Status: response.StatusCode, Method: http.MethodPatch, Url: uri, Error: err,}
	}

	jsonUnMarshallError:=json.Unmarshal(responseBody, typeVar)
	if jsonUnMarshallError != nil {
		return IntegrationStatus{ Status: response.StatusCode, Method: http.MethodPatch, Url: uri, Error: jsonUnMarshallError,}
	}

	return IntegrationStatus{
		Status:  response.StatusCode,
		Method:  http.MethodPatch,
		Url:     uri,
		Message: response.Status,
		Error:   nil,
	}



}