package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type ipresp struct {
	IP string `json:"ip"`
}

func TestEncodeURL(t *testing.T) {
	apiURL := "http://test.api/api?query=%s"
	payload := map[string]string{
		"the quick brown": fmt.Sprintf(apiURL, "the+quick+brown"),
	}
	for key, val := range payload {
		r := EncodeURL(apiURL, key)
		if r != val {
			t.Errorf("For %v : Expected %v, got %v", key, val, r)
		}
	}
}

func TestFetchURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"ip": "127.0.0.1"}`)
	}))
	defer ts.Close()
	var got ipresp
	FetchURL(ts.URL, &got)
	if !reflect.DeepEqual(ipresp{"127.0.0.1"}, got) {
		t.Fail()
	}
}
