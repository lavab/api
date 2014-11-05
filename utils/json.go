package utils

import (
	"encoding/json"
	"io"
	"log"
)

// ReadJSON reads a JSON string as a generic map string
func ReadJSON(r io.Reader) (map[string]interface{}, error) {
	decoder := json.NewDecoder(r)
	out := map[string]interface{}{}
	err := decoder.Decode(&out)
	if err != nil {
		return out, err
	}
	return out, nil
}
