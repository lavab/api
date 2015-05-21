package models

type Webhook struct {
	Resource

	Target  string `json:"target" gorethink:"target"`
	Type    string `json:"type" gorethink:"type"`
	Address string `json:"address" gorethink:"address"`
}
