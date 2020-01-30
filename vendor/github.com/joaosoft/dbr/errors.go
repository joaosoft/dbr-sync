package dbr

import "errors"

var (
	ErrorNumberOfConditionValues = errors.New("invalid number of condition values")
	ErrorEmptySet                = errors.New("empty set")
	ErrorInvalidPointer          = errors.New("the object isn't a pointer")
	ErrorDialectNotFound         = errors.New("dialect not found")
)
