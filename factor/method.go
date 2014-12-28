package factor

type Factor interface {
	Type() string
	Request(data string) (string, error)
	Verify(data []string, input string) (bool, error)
}
