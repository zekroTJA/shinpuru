package mody

import "errors"

var (
	ErrTypeMustBePointer = errors.New("v type must be a pointer")
	ErrFieldNotExistent  = errors.New("field not existent")
	ErrTypeMissmatch     = errors.New("field and value types are not matching")
)
