package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// EncodeURL encodes the url with the given query and apiURL
func EncodeURL(apiURL, query string) string {
	queryEnc := url.QueryEscape(query)
	return fmt.Sprintf(apiURL, queryEnc)
}

// FetchURL populates the given struct with the data at the given url
func FetchURL(url string, out interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Print("Something went wrong")
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(out)
	return err
}
