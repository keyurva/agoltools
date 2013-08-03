package agoltools

type DisplayError interface {
	DisplayError() string
}

type Error struct {
	Message string
	Code    int
	Cause   error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) DisplayError() string {
	return e.Message
}
