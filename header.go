package table

import "strings"

// A Header describes column data.
type Header []string

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
