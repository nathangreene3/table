package table2

// Row ...
type Row []interface{}

// NewRow ...
func NewRow(values ...interface{}) Row {
	r := make(Row, len(values))
	copy(r, values)
	return r
}

// Equal ...
func (r Row) Equal(row Row) bool {
	if len(r) != len(row) {
		return false
	}

	var i int
	for ; i < len(r) && r[i] == row[i]; i++ {
	}

	return i == len(r)
}
