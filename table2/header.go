package table2

// Header ...
type Header []string

// NewHeader ...
func NewHeader(ss ...string) Header {
	return append(make(Header, 0, len(ss)), ss...)
}

// Equal ...
func (h Header) Equal(hdr Header) bool {
	if len(h) != len(hdr) {
		return false
	}

	var i int
	for ; i < len(h) && h[i] == hdr[i]; i++ {
	}

	return i == len(h)
}
