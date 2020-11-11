package table

import (
	"strings"

	math "github.com/nathangreene3/math"
)

// A Column is a collection of the ith values in a body of rows.
type Column []interface{}

// NewColumn ...
func NewColumn(values ...interface{}) Column {
	c := make(Column, len(values))
	copy(c, values)
	return c
}

// Compare two rows by their strings representations.
func (c Column) Compare(col Column) int {
	var (
		a, b     = c.Strings(), col.Strings()
		maxIndex = math.MinInt(len(a), len(b))
	)

	for i := 0; i < maxIndex; i++ {
		switch {
		case a[i] < b[i]:
			return -1
		case b[i] < a[i]:
			return 1
		}
	}

	// Normally, m<n --> -1 and n<m --> 1, but here incomplete rows
	// should be pushed to the bottom of the table. For example,
	// [1 2 3] < [1 2] --> -1
	switch {
	case len(a) < len(b):
		return 1
	case len(b) < len(a):
		return -1
	default:
		return 0
	}
}

// CompareAt compares two rows by the ith value as strings.
func (c Column) CompareAt(col Column, i int) int {
	return strings.Compare(c.StringAt(i), col.StringAt(i))
}

// Copy a column.
func (c Column) Copy() Column {
	cpy := make(Column, len(c))
	copy(cpy, c)
	return cpy
}

// isEmpty determines if a column contains data or not.
func (c Column) isEmpty() bool {
	for i := 0; i < len(c); i++ {
		if c[i] != nil {
			switch baseTypeOf(c[i]) {
			case integerType, floatType:
				return false
			case stringType:
				if 0 < len(c[i].(string)) {
					return false
				}
			}
		}
	}

	return true
}

// String returns a string-representation of a column.
func (c Column) String() string {
	return strings.Join(c.Strings(), ",")
}

// StringAt returns the string-representation of the ith value.
func (c Column) StringAt(i int) string {
	return toString(c[i])
}

// Strings ...
func (c Column) Strings() []string {
	s := make([]string, 0, len(c))
	for i := 0; i < len(c); i++ {
		s = append(s, toString(c[i]))
	}

	return s
}

// Swap ...
func (c Column) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
