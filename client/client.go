package client

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/dchest/uniuri"
	"github.com/lavab/api/models"
	"github.com/lavab/api/routes"
	"github.com/lavab/sockjs-go-client"
)

type Client struct {
	sync.RWMutex

	Address  string
	SockJS   *sockjs.Client
	Headers  map[string]string
	Incoming map[string]chan *Response
	Handlers []func(ev *Event)
	Timeout  time.Duration

	subscription string
}

func New(address string, timeout time.Duration) (*Client, error) {
	sjs, err := sockjs.NewClient(address + "/ws")
	if err != nil {
		return nil, err
	}

	if timeout == 0 {
		timeout = time.Second * 30
	}

	client := &Client{
		Address:  address,
		SockJS:   sjs,
		Headers:  map[string]string{},
		Handlers: []func(ev *Event){},
		Incoming: map[string]chan *Response{},
		Timeout:  timeout,
	}

	go client.Loop()

	return client, nil
}

func (c *Client) Loop() {
	message := make(chan *Response)

	go func() {
		for {
			x := []string{}
			if err := c.SockJS.ReadMessage(&x); err != nil {
				log.Print(err)
				continue
			}

			var resp *Response
			if err := json.Unmarshal([]byte(x[0]), &resp); err != nil {
				log.Print(err)
				continue
			}

			message <- resp
		}
	}()

	select {
	case resp := <-message:
		if resp.Type == "response" {
			c.RLock()
			d, ok := c.Incoming[resp.ID]
			c.RUnlock()

			go func() {
				if ok {
					d <- resp
				}
			}()
		} else {
			// Run event handlers
			for _, handler := range c.Handlers {
				go func(handler func(ev *Event)) {
					handler(&Event{
						Type: resp.Type,
						ID:   resp.ID,
						Name: resp.Name,
					})
				}(handler)
			}
		}
	case <-c.SockJS.Reconnected:
		data, err := Encode(&Request{
			Type:  "subscribe",
			Token: c.subscription,
		})
		if err != nil {
			log.Print(err)
			return
		}

		if err := c.SockJS.WriteMessage(data); err != nil {
			log.Print(err)
			return
		}
	}
}

func (c *Client) Receive(id string) (*Response, error) {
	c.Lock()
	c.Incoming[id] = make(chan *Response)
	c.Unlock()

	select {
	case <-time.After(c.Timeout):
		return nil, errors.New("Request timeout")
	case data := <-c.Incoming[id]:
		return data, nil
	}
}

func (c *Client) Request(method, path string, headers map[string]string, body interface{}) ([]string, string, error) {
	if headers == nil {
		headers = map[string]string{}
	}

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
			return nil, "", err
		}

		req.Body = string(data)
	}

	d, err := Encode(req)
	if err != nil {
		return nil, "", err
	}

	return d, req.ID, nil
}

type Event struct {
	Type string `json:"event"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) Subscribe(token string, callback func(ev *Event)) error {
	c.Handlers = append(c.Handlers, callback)

	if token != "" {
		data, err := Encode(&Request{
			Type:  "subscribe",
			Token: token,
		})
		if err != nil {
			return err
		}

		if err := c.SockJS.WriteMessage(data); err != nil {
			return err
		}

		c.subscription = token
	}

	return nil
}

func (c *Client) CreateToken(req *routes.TokensCreateRequest) (*models.Token, error) {
	data, id, err := c.Request("POST", "/tokens", map[string]string{
		"Content-Type": "application/json;charset=utf-8",
	}, req)
	if err != nil {
		return nil, err
	}

	if err := c.SockJS.WriteMessage(data); err != nil {
		return nil, err
	}

	rcv, err := c.Receive(id)
	if err != nil {
		return nil, err
	}

	var resp *routes.TokensCreateResponse
	if err := json.Unmarshal([]byte(rcv.Body), &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp.Token, nil
}

func (c *Client) CreateEmail(req *routes.EmailsCreateRequest) ([]string, error) {
	data, id, err := c.Request("POST", "/emails", map[string]string{
		"Content-Type": "application/json;charset=utf-8",
	}, req)
	if err != nil {
		return nil, err
	}

	if err := c.SockJS.WriteMessage(data); err != nil {
		return nil, err
	}

	rcv, err := c.Receive(id)
	if err != nil {
		return nil, err
	}

	var resp *routes.EmailsCreateResponse
	if err := json.Unmarshal([]byte(rcv.Body), &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp.Created, nil
}

func (c *Client) GetKey(id string) (*models.Key, error) {
	data, id, err := c.Request("GET", "/keys/"+id, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.SockJS.WriteMessage(data); err != nil {
		return nil, err
	}

	rcv, err := c.Receive(id)
	if err != nil {
		return nil, err
	}

	var resp *routes.KeysGetResponse
	if err := json.Unmarshal([]byte(rcv.Body), &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp.Key, nil
}

func (c *Client) GetEmail(id string) (*models.Email, error) {
	data, id, err := c.Request("GET", "/emails/"+id, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.SockJS.WriteMessage(data); err != nil {
		return nil, err
	}

	rcv, err := c.Receive(id)
	if err != nil {
		return nil, err
	}

	var resp *routes.EmailsGetResponse
	if err := json.Unmarshal([]byte(rcv.Body), &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp.Email, nil
}
