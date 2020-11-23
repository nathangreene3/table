package table

import "time"

const (
	// Int ...
	Int Type = iota + 1

	// Flt ...
	Flt

	// Bool ...
	Bool

	// Time ...
	Time

	// Str ...
	Str
)

// Type ...
type Type byte

// Types ...
type Types []Type

// NewTypes ...
func NewTypes(ts ...Type) Types {
	return append(make(Types, 0, len(ts)), ts...)
}

// Copy ...
func (ts Types) Copy() Types {
	return append(make(Types, 0, len(ts)), ts...)
}

// Equal ...
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

// Parse ...
func Parse(x interface{}) Type {
	switch x.(type) {
	case int:
		return Int
	case float64:
		return Flt
	case bool:
		return Bool
	case time.Time:
		return Time
	case string:
		return Str
	default:
		return 0
	}
}
