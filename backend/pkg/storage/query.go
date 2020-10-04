package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// Query queries the storage
func Query(method, path string, reqbody, respbody interface{}) error {
	var reqbodystream io.Reader
	switch reqbody := reqbody.(type) {
	case io.Reader:
		reqbodystream = reqbody
	case string:
		reqbodystream = strings.NewReader(reqbody)
	case nil:
		// do nothing
	default:
		reqbodybytes, err := json.Marshal(reqbody)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reqbodystream = strings.NewReader(string(reqbodybytes))
	}

	url := viper.GetString("elastic.addr") + path
	req, err := http.NewRequest(method, url, reqbodystream)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	if reqbodystream != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respbytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf(
			"invalid response code: expected %d, got %d: %s",
			http.StatusOK,
			resp.StatusCode,
			string(respbytes),
		)
	}

	if respbody == nil {
		respbody = &map[string]interface{}{}
	}
	if err := json.NewDecoder(resp.Body).Decode(respbody); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}
	return nil
}

// Get request to storage
func Get(path string, respbody interface{}) error {
	return Query(http.MethodGet, path, nil, respbody)
}

// Post request to storage
func Post(path string, reqbody, respbody interface{}) error {
	return Query(http.MethodPost, path, reqbody, respbody)
}

// Put request to storage
func Put(path string, reqbody, respbody interface{}) error {
	return Query(http.MethodPut, path, reqbody, respbody)
}

// Delete request to storage
func Delete(path string, reqbody, respbody interface{}) error {
	return Query(http.MethodDelete, path, reqbody, respbody)
}
