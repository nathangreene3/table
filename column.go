package table

import (
	"fmt"
	"strings"
)

// A Column is a collection of the ith values in a body of rows.
type Column []interface{}

// Copy a column.
func (c Column) Copy() Column {
	cpy := make(Column, len(c))
	copy(cpy, c)
	return cpy
}

// String returns a string-representation of a column.
func (c Column) String() string {
	sb := strings.Builder{}
	for i := range c {
		sb.WriteString(fmt.Sprintf("%v", c[i]))
	}

	return sb.String()
}
