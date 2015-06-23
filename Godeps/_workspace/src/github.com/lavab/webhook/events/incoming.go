package events

type Incoming struct {
	Email   string `json:"email"`
	Account string `json:"account"`
}
