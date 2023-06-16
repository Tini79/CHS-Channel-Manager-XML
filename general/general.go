package general

import (
	"net/http"
	"strings"
)

func SetHeader(req *http.Request) *http.Request {
	req.Header.Set("Content-Type", "application/json")
	// TODO Token
	req.Header.Set("Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODc1MDM2MzAsInJlZnJlc2giOmZhbHNlLCJ1c2VyIjoiU1lTVEVNIn0.fqA0wHYfmtxhfD13M7zBXOxVan0OxeWY0elzAvGKTxk")
	return req
}

func GetBedTypeCode(InputString string) (string, string) {
	output := strings.Split(InputString, "#")
	return output[0], output[1]
}

func GenerateToken() {

}

var IP = "http://192.168.1.57"
var Port = "9000"
