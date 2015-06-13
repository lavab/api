package sockjs

type Connection interface {
	ReadJSON(interface{}) error
	WriteJSON(interface{}) error
	Close() error
}
