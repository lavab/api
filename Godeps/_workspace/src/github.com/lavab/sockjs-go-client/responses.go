package sockjs

type Info struct {
	WebSocket    bool     `json:"websocket"`
	CookieNeeded bool     `json:"cookie_needed"`
	Origins      []string `json:"origins"`
	Entropy      int      `json:"entropy"`
}
