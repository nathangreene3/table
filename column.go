package table

// Column ...
type Column []interface{}

// NewCol ...
func NewCol(values ...interface{}) Column {
	if 0 < len(values) {
		for i, f := 1, Parse(values[0]); i < len(values); i++ {
			if f != Parse(values[i]) {
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

// Type ...
func (c Column) Type() Type {
	var t Type
	if 0 < len(c) {
		t = Parse(c[0])
		for i := 1; i < len(c); i++ {
			if Parse(c[i]) != t {
				t = 0
				break
			}
		}
	}

	return t
}
