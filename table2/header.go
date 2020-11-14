package table2

import (
	"bytes"
	"encoding/json"
	"strings"
)

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

// MarshalJSON ...
func (h Header) MarshalJSON() ([]byte, error) {
	return []byte(h.JSON()), nil
}

// UnmarshalJSON ...
func (h Header) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		json.Marshal(nil)
		return &json.MarshalerError{}
	}

	bs := bytes.SplitN(b[1:len(b)-1], []byte(","), cap(h)-len(h))
	for i := len(h); i < len(bs); i++ {
		h = append(h, string(bs[i]))
	}

	return nil
}

// JSON ...
func (h Header) JSON() string {
	return "[" + strings.Join([]string(h), ",") + "]"
}
