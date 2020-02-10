package table

import (
	"strings"

	"github.com/nathangreene3/math"
)

// A Row is a single entry in a table.
type Row []interface{}

// NewRow converts various items to a row.
func NewRow(values ...interface{}) Row {
	r := make(Row, len(values))
	copy(r, values)
	return r
}

// RowFromBts ...
func RowFromBts(line []byte) Row {
	return RowFromStr(string(line))
}

// RowFromStr ...
func RowFromStr(line string) Row {
	var (
		ss = strings.Split(line, ",")
		r  = make(Row, 0, len(ss))
	)

	for _, s := range ss {
		r = append(r, parse(s))
	}

	return r
}

// Compare two rows by their strings representations.
func (r Row) Compare(row Row) int {
	var (
		a, b     = r.Strings(), row.Strings()
		m, n     = len(a), len(b)
		maxIndex = math.MinInt(m, n)
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
	case m < n:
		return 1
	case n < m:
		return -1
	default:
		return 0
	}
}

// CompareAt compares two rows by the ith value as strings.
func (r Row) CompareAt(row Row, i int) int {
	return strings.Compare(r.StringAt(i), row.StringAt(i))
}

// Copy a row.
func (r Row) Copy() Row {
	return NewRow(r...)
}

// isEmpty determines if a row contains data or not.
func (r Row) isEmpty() bool {
	for _, v := range r {
		if v == nil {
			continue
		}

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
	return strings.Join(r.Strings(), ",")
}

// StringAt returns the string-representation of the ith value.
func (r Row) StringAt(i int) string {
	return toString(r[i])
}

// Strings returns the row formatted as []string.
func (r Row) Strings() []string {
	s := make([]string, 0, len(r))
	for _, v := range r {
		s = append(s, toString(v))
	}

	return s
}

// Swap ...
func (r Row) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
