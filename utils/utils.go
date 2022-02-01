package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ParseBody is used to convert the JSON request into a structure that
// Go will be able to understand
func ParseBody(x interface{}, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, x); err != nil {
		return err
	}
	return nil
}
