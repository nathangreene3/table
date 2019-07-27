package table

import (
	"fmt"
	"strconv"
	"strings"
)

// A Row is a single entry in a table.
type Row []interface{}

// isEmpty determines if a row contains data or not.
func (r Row) isEmpty() bool {
	for _, v := range r {
		switch baseTypeOf(v) {
		case integerType:
			return false
		case floatType:
			return false
		}
	}

	return true
}

// Copy a row.
func (r Row) Copy() Row {
	cpy := make(Row, len(r))
	copy(cpy, r)
	return cpy
}

// String returns a string-representation of a row.
func (r Row) String() string {
	sb := strings.Builder{}
	for i := range r {
		sb.WriteString(fmt.Sprintf("%v", r[i]))
	}

	return sb.String()
}

// SliceString returns the row formatted as []string.
func (r Row) SliceString() []string {
	s := make([]string, 0, len(r))
	var x string
	for _, v := range r {
		switch baseTypeOf(v) {
		case integerType:
			s = append(s, strconv.Itoa(v.(int)))
		case floatType:
			x = strconv.FormatFloat(v.(float64), 'f', -1, 64)
			if strings.ContainsRune(x, '.') {
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
