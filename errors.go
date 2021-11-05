package table

const (
	// errDims indicates a slice does not have the same length as
	// another.
	errDims = "dimension mismatch"

	// errRange indicates an index is either too small or large to
	// access a value in an indexible object.
	errRange = "index out of range"

	// errTimeFmt indicates a time string failed to parse.
	errTimeFmt = "invalid time format"

	// errType indicates a type does not match another type.
	errType = "invalid type"

	// errVarCount indicates an unexpected number variadic arguments
	// were provided.
	errVarCount = "unexpected variadic argument provided"
)
