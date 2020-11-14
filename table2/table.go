package table2

import (
	"encoding/csv"
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

// AppendCol ...
func (t *Table) AppendCol(colName string, c Column) *Table {
	m, n := t.Dims()
	if m != len(c) {
		panic("dimension mismatch")
	}

	t.header = append(t.header, colName)
	if 0 < m {
		t.formats = append(t.formats, Fmt(c[0]))
		if m+len(t.body) <= cap(t.body) {
			t.body = append(t.body, make([]interface{}, m)...)
			for i := m - 1; 0 < i; i-- {
				t.body[i*n+i+n] = c[i]
				copy(t.body[i*n+i:i*n+i+n], t.body[i*n:i*n+n])
			}

			t.body[n] = c[0]
		} else {
			b := make(Body, 0, (m+1)*n)
			for i := 0; i < m; i++ {
				b = append(append(b, t.body[i*n:(i+1)*n]...), c[i])
			}

			t.body = b
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

// Export ...
func (t *Table) Export(w *csv.Writer) error {
	return w.WriteAll(t.Strings())
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

// Import ...
func Import(r *csv.Reader) (*Table, error) {
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return New(NewHeader()), nil
	}

	t := New(NewHeader(lines[0]...))
	for i := 1; i < len(lines); i++ {
		r := make(Row, 0, len(lines[i]))
		for j := 0; j < len(lines[i]); j++ {
			if n, err := strconv.ParseInt(lines[i][j], 10, strconv.IntSize); err == nil {
				r = append(r, int(n))
				continue
			}

			if f, err := strconv.ParseFloat(lines[i][j], strconv.IntSize); err == nil {
				r = append(r, f)
				continue
			}

			if b, err := strconv.ParseBool(lines[i][j]); err == nil {
				r = append(r, b)
				continue
			}

			r = append(r, lines[i][j])
		}

		t.Append(r)
	}

	return t, nil
}

// Insert ...
func (t *Table) Insert(i int, r Row) *Table {
	m, _ := t.Dims()
	return t.Append(r).Swap(i, m)
}

// InsertCol ...
func (t *Table) InsertCol(i int, colName string, c Column) *Table {
	_, n := t.Dims()
	return t.AppendCol(colName, c).SwapCols(i, n)
}

// JSON ...
func (t *Table) JSON() string {
	var (
		sb   strings.Builder
		m, n = t.Dims()
	)

	sb.Grow((m + 3) * n * 8)
	sb.WriteString("{\"header\":[")
	if 0 < n {
		sb.WriteString("\"" + strings.Join([]string(t.header), "\",\"") + "\"")
	}

	sb.WriteString("],\"formats\":[")
	if 0 < n {
		sb.WriteString("\"" + string(t.formats[0]) + "\"")
		for j := 1; j < n; j++ {
			sb.WriteString(",\"" + string(t.formats[j]) + "\"")
		}
	}

	sb.WriteString("],\"body\":[")
	if 0 < m && 0 < n {
		sb.WriteByte('[')
		switch t.formats[0] {
		case Int:
			sb.WriteString(strconv.FormatInt(int64(t.body[0].(int)), 10))
		case Flt:
			sb.WriteString(strconv.FormatFloat(float64(t.body[0].(float64)), 'f', -1, 64))
		case Bool:
			sb.WriteString("\"" + strconv.FormatBool(t.body[0].(bool)) + "\"")
		case Str:
			sb.WriteString("\"" + t.body[0].(string) + "\"")
		default:
		}

		for j := 1; j < n; j++ {
			switch t.formats[j] {
			case Int:
				sb.WriteString("," + strconv.FormatInt(int64(t.body[j].(int)), 10))
			case Flt:
				sb.WriteString("," + strconv.FormatFloat(t.body[j].(float64), 'f', -1, 64))
			case Bool:
				sb.WriteString(",\"" + strconv.FormatBool(t.body[j].(bool)) + "\"")
			case Str:
				sb.WriteString(",\"" + t.body[j].(string) + "\"")
			default:
			}
		}

		sb.WriteByte(']')
		for i := 1; i < m; i++ {
			sb.WriteString(",[")
			switch t.formats[0] {
			case Int:
				sb.WriteString(strconv.FormatInt(int64(t.body[i*n].(int)), 10))
			case Flt:
				sb.WriteString(strconv.FormatFloat(t.body[i*n].(float64), 'f', -1, 64))
			case Bool:
				sb.WriteString("\"" + strconv.FormatBool(t.body[i*n].(bool)) + "\"")
			case Str:
				sb.WriteString("\"" + t.body[i*n].(string) + "\"")
			default:
			}

			for j := 1; j < n; j++ {
				switch t.formats[j] {
				case Int:
					sb.WriteString("," + strconv.FormatInt(int64(t.body[i*n+j].(int)), 10))
				case Flt:
					sb.WriteString("," + strconv.FormatFloat(t.body[i*n+j].(float64), 'f', -1, 64))
				case Bool:
					sb.WriteString(",\"" + strconv.FormatBool(t.body[i*n+j].(bool)) + "\"")
				case Str:
					sb.WriteString(",\"" + t.body[i*n+j].(string) + "\"")
				default:
				}
			}

			sb.WriteString("]")
		}
	}

	sb.WriteString("]}")
	return sb.String()
}

// MarshalJSON ...
func (t *Table) MarshalJSON() ([]byte, error) {
	return []byte(t.JSON()), nil
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

// RemoveCol ...
func (t *Table) RemoveCol(j int) (string, Column) {
	var (
		m, n   = t.Dims()
		name   = t.header[j]
		column = make([]interface{}, 0, m)
	)

	t.header = append(t.header[:j], t.header[j+1:]...)
	t.formats = append(t.formats[:j], t.formats[j+1:]...)
	if 0 < m {
		for i := 0; i+1 < m; i++ {
			ij := i*n + j
			column = append(column, t.body[ij])
			copy(t.body[ij-i:ij-i+n-1], t.body[ij+1:ij+n])
		}

		ij := len(t.body) - n + j
		column = append(column, t.body[ij])
		t.body = append(t.body[:ij-m+1], t.body[ij+1:]...)
	}

	return name, column
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
func (t *Table) Strings() [][]string {
	var (
		m, n = t.Dims()
		ss   = make([][]string, 0, m+1)
		h    = make([]string, 0, n)
	)

	for j := 0; j < len(t.header); j++ {
		h = append(h, t.header[j])
	}

	ss = append(ss, h)
	for i := 0; i < m; i++ {
		r := make([]string, 0, n)
		for j := 0; j < n; j++ {
			switch v := t.body[i*n+j].(type) {
			case bool:
				r = append(r, strconv.FormatBool(v))
			case float64:
				r = append(r, strconv.FormatFloat(v, 'f', -1, 64))
			case int:
				r = append(r, strconv.Itoa(v))
			case string:
				r = append(r, v)
			default:
			}
		}

		ss = append(ss, r)
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

// SwapCols ...
func (t *Table) SwapCols(i, j int) *Table {
	t.header[i], t.header[j] = t.header[j], t.header[i]
	t.formats[i], t.formats[j] = t.formats[j], t.formats[i]

	_, n := t.Dims()
	for kn := 0; kn < len(t.body); kn += n {
		t.body[kn+i], t.body[kn+j] = t.body[kn+j], t.body[kn+i]
	}

	return t
}
