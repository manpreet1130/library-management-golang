package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

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

