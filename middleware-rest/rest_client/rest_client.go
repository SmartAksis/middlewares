package rest_client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

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


func Get(rest RestConsumer, url string, typeVar interface{}) error {
	uri:=rest.BaseEndpoint() + url
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jsonUnMarshallError:=json.Unmarshal(body, typeVar)
	if jsonUnMarshallError != nil {
		return jsonUnMarshallError
	}
	return nil
}