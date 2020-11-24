package table

import "errors"

var (
	// errDims indicates a slice does not have the same length as another.
	errDims = errors.New("dimension mismatch")

	// errRange indicates an index is either too small or large to access a
	// value in an indexible object.
	errRange = errors.New("index out of range")

	// errType indicates a type does not match another type.
	errType = errors.New("invalid type")
)
