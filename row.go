package table

import (
	"bytes"
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

	for i := 0; i < len(ss); i++ {
		r = append(r, parse(ss[i]))
	}

	return r
}

// Bytes ...
func (r Row) Bytes() []byte {
	if len(r) == 0 {
		return []byte{}
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	buf.Write(toBytes(r[0]))
	for i := 1; i < len(r); i++ {
		buf.WriteByte(',')
		buf.Write(toBytes(r[i]))
	}

	return buf.Bytes()
}

// Compare two rows by their strings representations.
func (r Row) Compare(row Row) int {
	var (
		a, b     = r.Strings(), row.Strings()
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

	// Normally, m < n --> -1 and n < m --> 1, but here incomplete rows
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
func (r Row) CompareAt(row Row, i int) int {
	return strings.Compare(r.StringAt(i), row.StringAt(i))
}

// Copy a row.
func (r Row) Copy() Row {
	return NewRow(r...)
}

// isEmpty determines if a row contains data or not.
func (r Row) isEmpty() bool {
	for i := 0; i < len(r); i++ {
		if r[i] != nil {
			switch baseTypeOf(r[i]) {
			case integerType, floatType:
				return false
			case stringType:
				if 0 < len(r[i].(string)) {
					return false
				}
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

// Strings returns the row updateBaseTypesAndWidthsted as []string.
func (r Row) Strings() []string {
	s := make([]string, 0, len(r))
	for i := 0; i < len(r); i++ {
		s = append(s, toString(r[i]))
	}

	return s
}

// Swap ...
func (r Row) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
