package utils

import (
	// "github.com/aurieh/ddg-ng/commandclient"
	"encoding/json"
	"net/http"
	"time"
)

// Client global http client for ddg
var Client = &http.Client{
	Timeout: time.Second * 2,
}

// GetJSON simplifies JSON decoding
func GetJSON(res *http.Response, target interface{}) error {
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(target)
}
