package table

import (
	"fmt"
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
