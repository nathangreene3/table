package table

// Row ...
type Row []interface{}

// NewRow ...
func NewRow(values ...interface{}) Row {
	return append(make(Row, 0, len(values)), values...)
}

// Copy ...
func (r Row) Copy() Row {
	return append(make(Row, 0, len(r)), r...)
}

// Equal ...
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

// Types ...
func (r Row) Types() Types {
	f := make(Types, 0, len(r))
	for i := 0; i < len(r); i++ {
		f = append(f, Parse(r[i]))
	}

	return f
}
