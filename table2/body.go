package table2

import (
	"encoding/json"
	"strconv"
)

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

// String ...
// func (b Body) String() string {
// 	var sb strings.Builder
// 	if 0 < len(b) {
// 		sb.Grow(len(b) * len(b[0]) * 8)
// 		for i := 0; i < len(b); i++ {

// 		}
// 	}

// 	return "[" + sb.String() + "]"
// }

// Strings ...
func (b Body) Strings() []string {
	ss := make([]string, 0, len(b))
	for i := 0; i < len(b); i++ {
		switch Fmt(b[i]) {
		case Bool:
			ss = append(ss, strconv.FormatBool(b[i].(bool)))
		case Flt:
			ss = append(ss, strconv.FormatFloat(b[i].(float64), 'f', -1, 64))
		case Int:
			ss = append(ss, strconv.Itoa(b[i].(int)))
		case Str:
			ss = append(ss, b[i].(string))
		default:
			panic("invalid format")
		}
	}

	return ss
}

// MarshalJSON ...
func (b Body) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Strings())
}
