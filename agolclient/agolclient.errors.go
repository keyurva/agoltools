package agolclient

import (
	"bytes"
	"fmt"
	"strings"
)

type Error struct {
	Message string
	Cause   error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) DisplayError() string {
	return e.Message
}

func DisplayError(message string, cause error) error {
	if r, ok := cause.(*RESTError); ok {
		return r
	}
	if e, ok := cause.(*Error); ok {
		return e
	}
	return &Error{Message: message, Cause: cause}
}

type RESTError struct {
	Code    int
	Message string
	Details []string
}

func (re *RESTError) Error() string {
	return fmt.Sprintf("%s, %d (%s)", re.Message, re.Code, re.Details)
}

func (re *RESTError) DisplayError() string {
	var buf bytes.Buffer
	buf.WriteString(re.Message)
	if re.Details != nil && len(re.Details) > 0 {
		buf.WriteString(fmt.Sprintf(" (%s)", strings.Join(re.Details, ", ")))
	}
	return buf.String()
}
