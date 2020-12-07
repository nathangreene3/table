package table

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

// Strings returns a list of strings in a header.
func (h Header) Strings() []string {
	return append(make([]string, 0, len(h)), h...)
}
