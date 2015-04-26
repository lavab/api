package client

import (
	"encoding/json"
	"log"

	"github.com/dchest/uniuri"
	"github.com/lavab/api/models"
	"github.com/lavab/api/routes"
	"github.com/lavab/sockjs-go-client"
)

type Client struct {
	Address string
	SockJS  *sockjs.Client
	Headers map[string]string
}

func New(address string) (*Client, error) {
	sjs, err := sockjs.NewClient(address + "/ws")
	if err != nil {
		return nil, err
	}

	return &Client{
		Address: address,
		SockJS:  sjs,
	}, nil
}

func (c *Client) Request(method, path string, headers map[string]string, body interface{}) ([]string, error) {
	if c.Headers != nil {
		for k, v := range c.Headers {
			headers[k] = v
		}
	}

	req := &Request{
		ID:      uniuri.New(),
		Type:    "request",
		Path:    path,
		Method:  method,
		Headers: headers,
	}

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		req.Body = string(data)
	}

	return Encode(req)
}

func (c *Client) CreateToken(req *routes.TokensCreateRequest) (*models.Token, error) {
	data, err := c.Request("POST", "/tokens", map[string]string{
		"Content-Type": "application/json;charset=utf-8",
	}, req)
	if err != nil {
		return nil, err
	}

	if err := c.SockJS.WriteMessage(data); err != nil {
		return nil, err
	}

	x := []string{}
	if err := c.SockJS.ReadMessage(&x); err != nil {
		return nil, err
	}

	var resp *routes.TokensCreateResponse
	if err := Decode(x, &resp); err != nil {
		return nil, err
	}

	return resp.Token, nil
}

func (c *Client) CreateEmail(req *routes.EmailsCreateRequest) ([]string, error) {
	data, err := c.Request("POST", "/emails", map[string]string{
		"Content-Type": "application/json;charset=utf-8",
	}, req)
	if err != nil {
		return nil, err
	}

	if err := c.SockJS.WriteMessage(data); err != nil {
		return nil, err
	}

	x := []string{}
	if err := c.SockJS.ReadMessage(&x); err != nil {
		return nil, err
	}

	var resp *routes.EmailsCreateResponse
	if err := Decode(x, &resp); err != nil {
		return nil, err
	}

	log.Printf("%+v", resp)

	return resp.Created, nil
}
