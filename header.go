package table

import "strings"

// Header is a list of column names.
type Header []string

// NewHeader returns a new header.
func NewHeader(ss ...string) Header {
	return append(make(Header, 0, len(ss)), ss...)
}

// Copy returns a copy of a header.
func (h Header) Copy() Header {
	return append(make(Header, 0, len(h)), h...)
}

// Equal determines if two headers are equal.
func (h Header) Equal(hdr Header) bool {
	if len(h) != len(hdr) {
		return false
	}

	for i := 0; i < len(h); i++ {
		if h[i] != hdr[i] {
			return false
		}
	}

	return true
}

// String ...
func (h Header) String() string {
	var sb strings.Builder
	if 0 < len(h) {
		n := len(h) + 3
		for i := 0; i < len(h); i++ {
			n += len(h[i])
		}

		sb.Grow(n)
		sb.WriteByte('[')
		for i := 0; i < len(h); i++ {
			sb.WriteString(" " + h[i])
		}

		sb.WriteString(" ]")
	} else {
		sb.WriteString("[]")
	}

	return sb.String()
}

// Strings returns a list of strings in a header.
func (h Header) Strings() []string {
	return append(make([]string, 0, len(h)), h...)
}
