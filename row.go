package table

// Row is a list of values.
type Row []interface{}

// NewRow returns a new row.
func NewRow(values ...interface{}) Row {
	return append(make(Row, 0, len(values)), values...)
}

// Copy returns a copy of a row.
func (r Row) Copy() Row {
	return append(make(Row, 0, len(r)), r...)
}

// Equal determines if two rows are equal.
func (r Row) Equal(row Row) bool {
	if len(r) != len(row) {
		return false
	}

	for i := 0; i < len(r); i++ {
		if r[i] != row[i] {
			return false
		}
	}

	return true
}

// Interfaces returns a list of values in a row.
func (r Row) Interfaces() []interface{} {
	return append(make([]interface{}, 0, len(r)), r...)
}

// Types returns a list of a row's types.
func (r Row) Types() Types {
	f := make(Types, 0, len(r))
	for i := 0; i < len(r); i++ {
		f = append(f, Parse(r[i]))
	}

	return f
}
