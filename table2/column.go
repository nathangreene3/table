package table2

// Column ...
type Column []interface{}

// NewCol ...
func NewCol(values ...interface{}) Column {
	return append(make(Column, 0, len(values)), values...)
}

// Copy ...
func (c Column) Copy() Column {
	return append(make(Column, 0, len(c)), c...)
}

// Equal ...
func (c Column) Equal(col Column) bool {
	if len(c) != len(col) {
		return false
	}

	for i := 0; i < len(c); i++ {
		if c[i] != col[i] {
			return false
		}
	}

	return true
}

// Fmts ...
func (c Column) Fmts() Formats {
	f := make(Formats, 0, len(c))
	for i := 0; i < len(c); i++ {
		f = append(f, Fmt(c[i]))
	}

	return f
}
