package table

import (
	math "github.com/nathangreene3/math"
)

// Body ...
type Body []Row

// NewBody ...
func NewBody(rows ...Row) Body {
	b := make(Body, 0, len(rows))
	for i := 0; i < len(rows); i++ {
		b = append(b, rows[i].Copy())
	}

	return b
}

// Compare ...
func (b Body) Compare(body Body) int {
	minIndex := math.MinInt(len(b), len(body))
	for i := 0; i < minIndex; i++ {
		if c := b[i].Compare(body[i]); c != 0 {
			return c
		}
	}

	switch {
	case len(b) < len(body):
		return -1
	case len(body) < len(b):
		return 1
	default:
		return 0
	}
}

// Copy ...
func (b Body) Copy() Body {
	return NewBody(b...)
}

// Swap ...
func (b Body) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// SwapCols ...
func (b Body) SwapCols(i, j int) {
	for i := 0; i < len(b); i++ {
		b[i].Swap(i, j)
	}
}
