package db

type ConnectionError struct {
	error
	WrongAuthKey bool
	Unreachable  bool
}

type NotFound struct {
	error
	Database bool
	Table    bool
	Item     bool
}
