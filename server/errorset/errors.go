package errorset

import (
	"encoding/json"
	"errors"
	"sync"
)

const (
	jsonType = "github.com/mwheatley3/ww/server/service.Error"
)

var (
	// ErrJSONInvalidUnmarshal indicates a json unmarshaling error
	ErrJSONInvalidUnmarshal = errors.New("Attempted to unmarshal a non service error into a service error")
)

func newError(code int, msg string) *Error {
	return &Error{msg, code}
}

// An Error is an ErrorSet error
type Error struct {
	message string
	code    int
}

// Error satisfies the builtin error interface
func (e *Error) Error() string { return e.message }

// Code returns the errors code
func (e *Error) Code() int { return e.code }

type jsonError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// MarshalJSON satisfies the encoding/json.Marshaler interface
func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonError{jsonType, e.message, e.code})
}

// UnmarshalJSON satisfies the encoding/json.Unmarshaler interface
func (e *Error) UnmarshalJSON(b []byte) error {
	var je jsonError

	if err := json.Unmarshal(b, &je); err != nil {
		return err
	}

	if je.Type != jsonType {
		return ErrJSONInvalidUnmarshal
	}

	e.code = je.Code
	e.message = je.Message

	return nil
}

// New returns a new ErrorSet
func New() *ErrorSet {
	return &ErrorSet{errs: []*Error{}}
}

// ErrorSet represents a set of distinct error values
type ErrorSet struct {
	l    sync.RWMutex
	errs []*Error
}

// New adds a new error to the error set
func (e *ErrorSet) New(msg string) *Error {
	e.l.Lock()
	defer e.l.Unlock()

	err := newError(len(e.errs), msg)
	e.errs = append(e.errs, err)

	return err
}

// FromCode returns an Error by error code
func (e *ErrorSet) FromCode(code int) *Error {
	e.l.RLock()
	defer e.l.RUnlock()

	if code >= len(e.errs) || code < 0 {
		return nil
	}

	return e.errs[code]
}

// FromJSON attempts to unmarshal an error from json. It pulls
// out the code and returns the err by code
func (e *ErrorSet) FromJSON(b []byte) (*Error, error) {
	var je jsonError

	if err := json.Unmarshal(b, &je); err != nil {
		return nil, err
	}

	if je.Type != jsonType {
		return nil, ErrJSONInvalidUnmarshal
	}

	return e.FromCode(je.Code), nil
}

// Member returns true if the error set is a member
// of the set
func (e *ErrorSet) Member(err error) bool {
	err2, ok := err.(*Error)

	if !ok {
		return false
	}

	e.l.RLock()
	defer e.l.RUnlock()

	return e.errs[err2.Code()] == err2
}
