package table2

// Column ...
type Column []interface{}

// NewCol ...
func NewCol(values ...interface{}) Column {
	if 0 < len(values) {
		for i, f := 1, Fmt(values[0]); i < len(values); i++ {
			if f != Fmt(values[i]) {
				panic("invalid format")
			}
		}
	}

	return append(make(Column, 0, len(values)), values...)
}

// Copy ...
func (c Column) Copy() Column {
	return append(make(Column, 0, len(c)), c...)
}

// Equal ...
func (c Column) Equal(col Column) bool {
	n := len(c)
	if n != len(col) {
		return false
	}

	for i := 0; i < n; i++ {
		if c[i] != col[i] {
			return false
		}
	}

	return true
}

// Fmt ...
func (c Column) Fmt() Format {
	var f Format
	if 0 < len(c) {
		f = Fmt(c[0])
	}

	return f
}
