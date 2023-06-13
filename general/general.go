package general

import (
	"net/http"
)

func GetHeader(req *http.Request) *http.Request {
	req.Header.Set("Content-Type", "application/json")
	// TODO Token
	req.Header.Set("Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODc1MDM2MzAsInJlZnJlc2giOmZhbHNlLCJ1c2VyIjoiU1lTVEVNIn0.fqA0wHYfmtxhfD13M7zBXOxVan0OxeWY0elzAvGKTxk")

	return req
}

var IP = "http://192.168.1.64:9000/"