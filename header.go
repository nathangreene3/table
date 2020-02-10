package table

import (
	"strings"

	math "github.com/nathangreene3/math"
)

// A Header describes column data.
type Header []string

// NewHeader converts a list of column names to a header.
func NewHeader(colNames ...string) Header {
	h := make(Header, len(colNames))
	copy(h, colNames)
	return h
}

// HeaderFromBts ...
func HeaderFromBts(colNames ...[]byte) Header {
	h := make(Header, 0, len(colNames))
	for _, name := range colNames {
		h = append(h, string(name))
	}

	return h
}

// Compare ...
func (h Header) Compare(header Header) int {
	var (
		m, n     = len(h), len(header)
		minIndex = math.MinInt(m, n)
	)

	for i := 0; i < minIndex; i++ {
		if c := strings.Compare(h[i], header[i]); c != 0 {
			return c
		}
	}

	switch {
	case m < n:
		return -1
	case n < m:
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
	for _, v := range h {
		if len(v) != 0 {
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
