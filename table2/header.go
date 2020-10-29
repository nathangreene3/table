package table2

// Header ...
type Header []string

// NewHeader ...
func NewHeader(ss ...string) Header {
	return append(make(Header, 0, len(ss)), ss...)
}

// Copy ...
func (h Header) Copy() Header {
	return append(make(Header, 0, len(h)), h...)
}

// Equal ...
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
