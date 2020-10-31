package table2

import (
	"strconv"
	"strings"
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

// AppendCol ...TODO
// func (t *Table) AppendCol(colName string, c Column) *Table {
// 	m, n := t.Dims()
// 	if m != len(c) {
// 		panic("dimension mismatch")
// 	}

// 	t.header = append(t.header, colName)
// 	if 0 < m {
// 		t.formats = append(t.formats, Fmt(c[0]))
// 		if len(c)+len(t.body) <= cap(t.body) {
// 			t.body = append(t.body, make([]interface{}, m)...)
// 			for i := m - 1; 0 <= i; i-- {
// 				t.body[i*n+]
// 				copy(t.body[i*] ,t.body[i*n : i*(n+1)])
// 			}
// 		} else {

// 		}
// 	}

// 	return t
// }

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
func (t *Table) ColFmt(j int) Format {
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

// Fmts ...
func (t *Table) Fmts() Formats {
	return t.formats.Copy()
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

// Header ...
func (t *Table) Header() Header {
	return t.header.Copy()
}

// Insert ...
func (t *Table) Insert(i int, r Row) *Table {
	m, _ := t.Dims()
	return t.Append(r).Swap(i, m)
}

// InsertCol ...
func (t *Table) InsertCol(i int, c Column) *Table {
	return t
}

// Remove ...
func (t *Table) Remove(i int) Row {
	var (
		_, n = t.Dims()
		r    = NewRow(t.body[i*n : (i+1)*n]...)
	)

	t.body = append(t.body[:i*n], t.body[(i+1)*n:]...)
	if len(t.body) < n {
		t.formats = t.formats[:0]
	}

	return r
}

// Row ...
func (t *Table) Row(i int) Row {
	_, n := t.Dims()
	return NewRow(t.body[i*n : (i+1)*n]...)
}

// Rows ...
func (t *Table) Rows() []Row {
	var (
		m, n = t.Dims()
		rs   = make([]Row, 0, m)
	)

	for i := 0; i < m; i++ {
		rs = append(rs, NewRow(t.body[i*n:(i+1)*n]))
	}

	return rs
}

// Set ...
func (t *Table) Set(i, j int, v interface{}) *Table {
	if Fmt(v) != t.formats[j] {
		panic("invalid format")
	}

	_, n := t.Dims()
	t.body[i*n+j] = v
	return t
}

// Sort ...
func (t *Table) Sort(j int) *Table {
	// TODO
	return t
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

	h := "+" + strings.Join(hs, "+") + "+"

	var sb strings.Builder
	sb.WriteString(h + "\n|")
	for j := 0; j < n; j++ {
		switch t.formats[j] {
		case Flt, Int:
			sb.WriteString(strings.Repeat(" ", ws[j]-len(t.header[j])+1) + t.header[j] + " |")
		case Bool, Str:
			sb.WriteString(" " + t.header[j] + strings.Repeat(" ", ws[j]-len(t.header[j])+1) + "|")
		default:
		}
	}

	sb.WriteString("\n" + h)
	for i := 0; i < len(b); i += n {
		sb.WriteString("\n|")
		for j, k := 0, i; j < n && k < len(b); j, k = j+1, k+1 {
			switch t.formats[j] {
			case Flt, Int:
				// Right-aligned
				sb.WriteString(strings.Repeat(" ", ws[j]-len(b[k])+1) + b[k] + " |")
			case Bool, Str:
				// Left-aligned
				sb.WriteString(" " + b[k] + strings.Repeat(" ", ws[j]-len(b[k])+1) + "|")
			}
		}
	}

	sb.WriteString("\n" + h)
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

// Swap ...
func (t *Table) Swap(i, j int) *Table {
	_, n := t.Dims()
	for k := 0; k < n; k++ {
		ik, jk := i*n+k, j*n+k
		t.body[ik], t.body[jk] = t.body[jk], t.body[ik]
	}

	return t
}
