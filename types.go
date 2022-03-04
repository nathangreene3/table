package table

// TODO: Add support for complex numbers and money.

const (
	// Inv corresponds to any invalid type.
	Inv Type = iota

	// Int corresponds to type int.
	Int

	// Flt corresponds to type float64.
	Flt

	// Bool corresponds to type bool.
	Bool

	// Time corresponds to type time.Time.
	Time

	// Str corresponds to type string.
	Str
)

// Type corresponds to a basic type.
type Type byte

// Types is a list of types.
type Types []Type

// NewTypes returns a new list of types.
func NewTypes(ts ...Type) Types {
	return append(make(Types, 0, len(ts)), ts...)
}

// Copy returns a copy of a list of types.
func (ts Types) Copy() Types {
	return append(make(Types, 0, len(ts)), ts...)
}

// Equal determines if two lists of types are equal.
func (ts Types) Equal(types Types) bool {
	n := len(ts)
	if n != len(types) {
		return false
	}

	for i := 0; i < n; i++ {
		if ts[i] != types[i] {
			return false
		}
	}

	return true
}

// Parse returns the type of a given value.
func Parse(x interface{}) Type {
	switch x.(type) {
	case int:
		return Int
	case float64:
		return Flt
	case bool:
		return Bool
	case FTime:
		return Time
	case string:
		return Str
	default:
		return Inv
	}
}
