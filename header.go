package table

import (
	"strings"

	math "github.com/nathangreene3/math"
)

// A Header describes column data.
type Header []string

// NewHeader converts a list of column names to a header.
func NewHeader(colNames ...string) Header {
	r := strings.NewReplacer(",", "", "\n", "")
	h := make(Header, 0, len(colNames))
	for i := 0; i < len(colNames); i++ {
		h = append(h, r.Replace(colNames[i]))
	}

	return h
}

// HeaderFromBts ...
func HeaderFromBts(colNames ...[]byte) Header {
	r := strings.NewReplacer(",", "", "\n", "")
	h := make(Header, 0, len(colNames))
	for i := 0; i < len(colNames); i++ {
		h = append(h, r.Replace(string(colNames[i])))
	}

	return h
}

// Bytes ...
func (h Header) Bytes() []byte {
	return []byte(strings.Join(h, ","))
}

// Compare ...
func (h Header) Compare(header Header) int {
	minIndex := math.MinInt(len(h), len(header))
	for i := 0; i < minIndex; i++ {
		if c := strings.Compare(h[i], header[i]); c != 0 {
			return c
		}
	}

	switch {
	case len(h) < len(header):
		return -1
	case len(header) < len(h):
		return 1
	default:
		return 0
	}
}

// CompareAt ...
func (h Header) CompareAt(header Header, i int) int {
	return strings.Compare(h[i], header[i])
}

// Copy a header.
func (h Header) Copy() Header {
	return NewHeader(h...)
}

// isEmpty ...
func (h Header) isEmpty() bool {
	for i := 0; i < len(h); i++ {
		if 0 < len(h[i]) {
			return false
		}
	}

	return true
}

// String returns a string-representation of a header.
func (h Header) String() string {
	return strings.Join(h, ",")
}

// Swap ...
func (h Header) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
