package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Query queries the storage
func Query(method, path string, reqbody, respbody interface{}) error {
	var reqbodystream io.Reader
	if reqbody != nil {
		reqbodybytes, err := json.Marshal(reqbody)
		if err != nil {
			return errors.Wrap(err, "marshal request body")
		}
		reqbodystream = strings.NewReader(string(reqbodybytes))
	}

	url := viper.GetString("elastic.addr") + path
	req, err := http.NewRequest(method, url, reqbodystream)
	if err != nil {
		return errors.Wrap(err, "new request")
	}
	if reqbodystream != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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
		return errors.Wrap(err, "unmarshal response")
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

// Delete request to storage
func Delete(path string, reqbody, respbody interface{}) error {
	return Query(http.MethodDelete, path, reqbody, respbody)
}
