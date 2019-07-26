package table

import "strconv"

// baseType indicates the basic underlying type of a interface{} value.
type baseType byte

const (
	stringType baseType = iota
	floatType
	integerType
)

// toBaseType converts a string to the first type it converts to successfully. Preference is given as int, float64, string.
func toBaseType(s string) interface{} {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}

	if x, err := strconv.ParseFloat(s, 64); err == nil {
		return x
	}

	return s
}

// baseTypeOf returns format corresponding to the underlying type of x.
func baseTypeOf(x interface{}) baseType {
	switch x.(type) {
	case int:
		return integerType
	case float64:
		return floatType
	case string:
		return stringType
	default:
		panic("unknown type")
	}
}
