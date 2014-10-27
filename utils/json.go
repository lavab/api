package utils

import (
	"encoding/json"
	"io"
	"log"
)

// MakeJSON is a wrapper over json.Marshal that returns an error message if an error occured
func MakeJSON(data map[string]interface{}) string {
	res, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("Error marshalling the response body.", err)
		return "{\"status\":500,\"message\":\"Error occured while marshalling the response body\"}"
	}
	return string(res)
}

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
