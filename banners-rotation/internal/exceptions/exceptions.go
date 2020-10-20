package exceptions

import "errors"

var (
	ErrObjectNotExist = errors.New("object doesn't exist")
	ErrNoChanges      = errors.New("zero changes for operation")
)
