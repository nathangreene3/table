package table2

import "strconv"

// Body ...
type Body []interface{}

// Strings ...
func (b Body) Strings() []string {
	ss := make([]string, 0, len(b))
	for i := 0; i < len(b); i++ {
		switch v := b[i].(type) {
		case float64:
			ss = append(ss, strconv.FormatFloat(v, 'f', -1, 64))
		case int:
			ss = append(ss, strconv.Itoa(v))
		case string:
			ss = append(ss, v)
		default:
		}
	}

	return ss
}

// Equal ...
func (b Body) Equal(bdy Body) bool {
	if len(b) != len(bdy) {
		return false
	}

	var i int
	for ; i < len(b) && b[i] == bdy[i]; i++ {
	}

	return i == len(b)
}
