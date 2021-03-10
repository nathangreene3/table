package table

import (
	"strconv"
	"strings"
)

// Body is a list of mn values.
type Body []interface{}

// NewBody returns a new body of values.
func NewBody(values ...interface{}) Body {
	return append(make(Body, 0, len(values)), values...)
}

// Copy returns a copy of a body.
func (b Body) Copy() Body {
	return append(make(Body, 0, len(b)), b...)
}

// Equal determines if two bodies are equal.
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

// String returns a body formatted as [ v0 v1 ... ]
func (b Body) String() string {
	var sb strings.Builder
	if 0 < len(b) {
		sb.Grow(256)
		sb.WriteByte('[')
		for i := 0; i < len(b); i++ {
			switch ParseType(b[i]) {
			case Int:
				sb.WriteString(" " + strconv.Itoa(b[i].(int)))
			case Flt:
				if v := b[i].(float64); v == float64(int(v)) {
					sb.WriteString(" " + strconv.FormatFloat(v, 'f', 1, 64)) // Forces f.0 when value is an integer
				} else {
					sb.WriteString(" " + strconv.FormatFloat(v, 'f', -1, 64))
				}
			case Bool:
				sb.WriteString(" " + strconv.FormatBool(b[i].(bool)))
			case Time:
				sb.WriteString(" " + b[i].(FTime).String())
			case Str:
				sb.WriteString(" " + b[i].(string))
			default:
				panic(errType.Error())
			}
		}

		sb.WriteString(" ]")
	} else {
		sb.WriteString("[]")
	}

	return sb.String()
}

// Strings returns a list of strings in which each string is the string value converted to string by parsing the type. If a value does not parse, the empty string will be inserted.
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
			panic(errType.Error())
		}
	}

	return ss
}

// Types returns a list of types of each value.
func (b Body) Types() Types {
	ts := make(Types, 0, len(b))
	for i := 0; i < len(b); i++ {
		ts = append(ts, ParseType(b[i]))
	}

	return ts
}
