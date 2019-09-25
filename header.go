package table

import "strings"

// A Header describes column data.
type Header []string

// NewHeader converts a list of column names to a header.
func NewHeader(colNames ...string) Header {
	h := make(Header, 0, len(colNames))
	for _, c := range colNames {
		h = append(h, c)
	}

	return h
}

// Copy a header.
func (h Header) Copy() Header {
	cpy := make(Header, 0, len(h))
	copy(cpy, h)
	return cpy
}

// String returns a string-representation of a header.
func (h Header) String() string {
	return strings.Join(h, " ")
}
