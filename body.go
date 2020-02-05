package table

import (
	math "github.com/nathangreene3/math"
)

// Body ...
type Body []Row

// NewBody ...
func NewBody(rows ...Row) Body {
	b := make(Body, 0, len(rows))
	for _, r := range rows {
		b = append(b, r.Copy())
	}

	return b
}

// Compare ...
func (b Body) Compare(body Body) int {
	var (
		m, n     = len(b), len(body)
		minIndex = math.MinInt(m, n)
	)

	for i := 0; i < minIndex; i++ {
		if c := b[i].Compare(body[i]); c != 0 {
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
	for _, r := range b {
		r.Swap(i, j)
	}
}
