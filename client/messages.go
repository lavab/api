package client

import (
	"encoding/json"
)

type Request struct {
	Type    string            `json:"type,omitempty"`
	ID      string            `json:"id,omitempty"`
	Method  string            `json:"method,omitempty"`
	Path    string            `json:"path,omitempty"`
	Body    string            `json:"body,omitempty"`
	Token   string            `json:"token,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type Response struct {
	Type    string              `json:"type,omitempty"`
	ID      string              `json:"id,omitempty"`
	Body    string              `json:"body,omitempty"`
	Name    string              `json:"name,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
}

func Encode(r interface{}) ([]string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return []string{string(data)}, nil
}

func Decode(x []string, r interface{}) error {
	var v *Response
	if err := json.Unmarshal([]byte(x[0]), &v); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(v.Body), &r); err != nil {
		return err
	}

	return nil
}
