package table

import (
	"fmt"
	"strconv"
	"strings"
)

// A Row is a single entry in a table.
type Row []interface{}

// NewRow converts various items to a row.
func NewRow(values ...interface{}) Row {
	r := make(Row, 0, len(values))
	for _, v := range values {
		r = append(r, v)
	}

	return r
}

// Copy a row.
func (r Row) Copy() Row {
	cpy := make(Row, len(r))
	copy(cpy, r)
	return cpy
}

// isEmpty determines if a row contains data or not.
func (r Row) isEmpty() bool {
	for _, v := range r {
		switch baseTypeOf(v) {
		case integerType, floatType:
			return false
		case stringType:
			if 0 < len(v.(string)) {
				return false
			}
		}
	}

	return true
}

// String returns a string-representation of a row.
func (r Row) String() string {
	var sb strings.Builder
	for _, v := range r {
		sb.WriteString(fmt.Sprintf("%v", v))
	}

	return sb.String()
}

// Strings returns the row formatted as []string.
func (r Row) Strings() []string {
	s := make([]string, 0, len(r))
	for _, v := range r {
		switch baseTypeOf(v) {
		case integerType:
			s = append(s, strconv.Itoa(v.(int)))
		case floatType:
			if x := strconv.FormatFloat(v.(float64), 'f', -1, 64); strings.ContainsRune(x, '.') {
				s = append(s, x)
			} else {
				s = append(s, x+".0")
			}
		case stringType:
			s = append(s, v.(string))
		default:
			panic("uknown base type")
		}
	}

	return s
}
