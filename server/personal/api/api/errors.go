package api

import (
	"github.com/mwheatley3/ww/server/errorset"
)

// Errs is the sms schema error set
var Errs = errorset.New()

// schema errors
var (
	ErrInvalidModelID   = Errs.New("No such model")
	ErrInvalidIntentID  = Errs.New("No such intent")
	ErrInvalidExampleID = Errs.New("No such intent example")

	ErrInvalidUserID   = Errs.New("Invalid user id")
	ErrInvalidEmail    = Errs.New("Invalid email")
	ErrInvalidPassword = Errs.New("Invalid password")
	ErrDuplicateUser   = Errs.New("User with that email already exists")

	ErrUnknownError = Errs.New("Unknown working wheatley error")
)
