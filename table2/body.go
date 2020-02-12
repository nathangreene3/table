package table2

// Body ...
type Body []Row

// NewBody ...
func NewBody(rows ...Row) Body {
	b := make(Body, 0, len(rows))
	for i := range rows {
		b = append(b, rows[i].Copy())
	}

	return b
}

// Copy ...
func (b Body) Copy() Body {
	return NewBody(b...)
}
