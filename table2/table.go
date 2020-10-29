package table2

import (
	"strconv"
	"strings"
)

const (
	// Flt ...
	Flt = "float64"

	// Int ...
	Int = "int"

	// Str ...
	Str = "string"
)

// Table ...
type Table struct {
	header  Header
	formats Formats
	body    Body
}

// New ...
func New(h Header, r ...Row) *Table {
	t := Table{
		header:  append(make(Header, 0, len(h)), h...),
		formats: make(Formats, 0, len(h)),
		body:    make(Body, 0, len(r)*len(h)),
	}

	return t.Append(r...)
}

// Append ...
func (t *Table) Append(r ...Row) *Table {
	if 0 < len(r) {
		var i int
		if len(t.body) == 0 {
			if len(t.header) != len(r[0]) {
				panic("dimension mismatch")
			}

			for j := 0; j < len(r[0]); j++ {
				t.formats = append(t.formats, Fmt(r[0][j]))
				t.body = append(t.body, r[0][j])
			}

			i++
		}

		for ; i < len(r); i++ {
			if len(t.header) != len(r[i]) {
				panic("dimension mismatch")
			}

			for j := 0; j < len(r[i]); j++ {
				if Fmt(r[i][j]) != t.formats[j] {
					panic("invalid format")
				}

				t.body = append(t.body, r[i][j])
			}
		}
	}

	return t
}

// Col ...
func (t *Table) Col(j int) []interface{} {
	m, n := t.Dims()
	if n <= j {
		panic("index out of range")
	}

	c := make([]interface{}, 0, m)
	for ; j < len(t.body); j += n {
		c = append(c, t.body[j])
	}

	return c
}

// ColFmt ...
func (t *Table) ColFmt(j int) string {
	return t.formats[j]
}

// ColFlts ...
func (t *Table) ColFlts(j int) []float64 {
	m, n := t.Dims()
	if n <= j {
		panic("index out of range")
	}

	c := make([]float64, 0, m)
	for ; j < len(t.body); j += n {
		c = append(c, t.body[j].(float64))
	}

	return c
}

// ColInts ...
func (t *Table) ColInts(j int) []int {
	m, n := t.Dims()
	if n <= j {
		panic("index out of range")
	}

	c := make([]int, 0, m)
	for ; j < len(t.body); j += n {
		c = append(c, t.body[j].(int))
	}

	return c
}

// ColStrs ...
func (t *Table) ColStrs(j int) []string {
	m, n := t.Dims()
	if n <= j {
		panic("index out of range")
	}

	c := make([]string, 0, m)
	for ; j < len(t.body); j += n {
		c = append(c, t.body[j].(string))
	}

	return c
}

// Copy ...
func (t *Table) Copy() *Table {
	cpy := Table{
		header:  append(make(Header, 0, len(t.header)), t.header...),
		formats: append(make(Formats, 0, len(t.formats)), t.formats...),
		body:    append(make(Body, 0, len(t.body)), t.body...),
	}

	return &cpy
}

// Dims ...
func (t *Table) Dims() (int, int) {
	m, n := len(t.body), len(t.header)
	if 0 < n {
		m /= n
	}

	return m, n
}

// Equal ...
func (t *Table) Equal(tbl *Table) bool {
	return t.header.Equal(tbl.header) && t.formats.Equal(tbl.formats) && t.body.Equal(tbl.body)
}

// Fmt ...
func (t *Table) Fmt() Formats {
	format := make(Formats, 0, len(t.header))
	for j := 0; j < len(t.header); j++ {
		format = append(format, t.formats[j])
	}

	return format
}

// Fmt ...
func Fmt(x interface{}) string {
	switch x.(type) {
	case float64:
		return Flt
	case int:
		return Int
	case string:
		return Str
	default:
		return ""
	}
}

// Get ...
func (t *Table) Get(i, j int) interface{} {
	return t.body[i*len(t.header)+j]
}

// GetFlt ...
func (t *Table) GetFlt(i, j int) float64 {
	return t.body[i*len(t.header)+j].(float64)
}

// GetInt ...
func (t *Table) GetInt(i, j int) int {
	return t.body[i*len(t.header)+j].(int)
}

// GetStr ...
func (t *Table) GetStr(i, j int) string {
	return t.body[i*len(t.header)+j].(string)
}

// Remove ...
func (t *Table) Remove(i int) Row {
	var (
		_, n = t.Dims()
		r    = NewRow(t.body[i*n : (i+1)*n]...)
	)

	t.body = append(t.body[:i*n], t.body[(i+1)*n:]...)
	return r
}

// Row ...
func (t *Table) Row(i int) Row {
	_, n := t.Dims()
	return NewRow(t.body[i*n : (i+1)*n]...)
}

// String ...
func (t *Table) String() string {
	var (
		n  = len(t.header)
		ws = make([]int, n)
	)

	for j := 0; j < n; j++ {
		if ws[j] < len(t.header[j]) {
			ws[j] = len(t.header[j])
		}
	}

	b := t.body.Strings()
	for i := 0; i < len(b); i++ {
		if ws[i%n] < len(b[i]) {
			ws[i%n] = len(b[i])
		}
	}

	hs := make([]string, 0, n)
	for j := 0; j < n; j++ {
		hs = append(hs, strings.Repeat("-", ws[j]+2))
	}

	h := "+" + strings.Join(hs, "+") + "+\n"

	var sb strings.Builder
	sb.WriteString(h + "|")
	for j := 0; j < n; j++ {
		switch t.formats[j] {
		case Flt, Int:
			sb.WriteString(strings.Repeat(" ", ws[j]-len(t.header[j])+1) + t.header[j] + " |")
		case Str:
			sb.WriteString(" " + t.header[j] + strings.Repeat(" ", ws[j]-len(t.header[j])+1) + "|")
		default:
		}
	}

	sb.WriteString("\n" + h)
	for i := 0; i < len(b); i += n {
		sb.WriteByte('|')
		for j, k := 0, i; j < n && k < len(b); j, k = j+1, k+1 {
			switch t.formats[j] {
			case Flt, Int:
				sb.WriteString(strings.Repeat(" ", ws[j]-len(b[k])+1) + b[k] + " |")
			case Str:
				sb.WriteString(" " + b[k] + strings.Repeat(" ", ws[j]-len(b[k])+1) + "|")
			default:
			}
		}

		sb.WriteByte('\n')
	}

	sb.WriteString(h)
	return sb.String()
}

// Strings ...
func (t *Table) Strings() []string {
	ss := make([]string, 0, len(t.header)+len(t.body))
	for j := 0; j < len(t.header); j++ {
		ss = append(ss, t.header[j])
	}

	for i := 0; i < len(t.body); i++ {
		switch v := t.body[i].(type) {
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
