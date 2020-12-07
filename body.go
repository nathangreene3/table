package table

import "strconv"

// Body ...
type Body []interface{}

// NewBody ...
func NewBody(values ...interface{}) Body {
	return append(make(Body, 0, len(values)), values...)
}

// Copy ...
func (b Body) Copy() Body {
	return append(make(Body, 0, len(b)), b...)
}

// Equal ...
func (b Body) Equal(bdy Body) bool {
	if len(b) != len(bdy) {
		return false
	}

	for i := 0; i < len(b); i++ {
		if b[i] != bdy[i] {
			return false
		}
	}

	return true
}

// Strings ...
func (b Body) Strings() []string {
	ss := make([]string, 0, len(b))
	for i := 0; i < len(b); i++ {
		switch ParseType(b[i]) {
		case Int:
			ss = append(ss, strconv.Itoa(b[i].(int)))
		case Flt:
			if v := b[i].(float64); v == float64(int(v)) {
				ss = append(ss, strconv.FormatFloat(v, 'f', 1, 64)) // Forces f.0 when value is an integer
			} else {
				ss = append(ss, strconv.FormatFloat(v, 'f', -1, 64))
			}
		case Bool:
			ss = append(ss, strconv.FormatBool(b[i].(bool)))
		case Time:
			ss = append(ss, b[i].(FTime).String())
		case Str:
			ss = append(ss, b[i].(string))
		default:
			ss = append(ss, "")
		}
	}

	return ss
}

// Types ...
func (b Body) Types() Types {
	ts := make(Types, 0, len(b))
	for i := 0; i < len(b); i++ {
		ts = append(ts, ParseType(b[i]))
	}

	return ts
}
