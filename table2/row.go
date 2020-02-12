package table2

// Row ...
type Row []Cell

// NewRow ...
func NewRow(values ...interface{}) Row {
	r := make(Row, 0, len(values))
	for i := range values {
		r = append(r, NewCell(values[i]))
	}

	return r
}

// Copy ...
func (r Row) Copy() Row {
	cpy := make(Row, len(r))
	copy(cpy, r)
	return cpy
}
