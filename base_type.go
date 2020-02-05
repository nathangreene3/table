package table

import (
	"strconv"
	"strings"
)

// baseType indicates the basic underlying type of a interface{} value.
type baseType byte

const (
	unknownType baseType = iota // 0
	stringType                  // 1
	floatType                   // 2
	integerType                 // 3
)

// parse converts a string to the first type it converts to successfully.
// Preference is given as int, float64, string.
func parse(s string) interface{} {
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
		return unknownType
	}
}

// toString converts an interface{} to a string.
func toString(x interface{}) string {
	switch baseTypeOf(x) {
	case integerType:
		return strconv.Itoa(x.(int))
	case floatType:
		f := strconv.FormatFloat(x.(float64), 'f', -1, 64)
		if strings.ContainsRune(f, '.') {
			return f
		}
		return f + ".0"
	case stringType:
		return x.(string)
	default:
		panic("unknown base type")
	}
}
