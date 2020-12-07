package table

import (
	"encoding/csv"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// A Table holds tabular data.
type Table struct {
	header Header
	types  Types
	body   Body
}

// ---------------------------------------------------------------------------
// Constructors
// ---------------------------------------------------------------------------

// New returns a new table.
func New(h Header, r ...Row) *Table {
	t := Table{
		header: append(make(Header, 0, len(h)), h...),
		types:  make(Types, 0, len(h)),
		body:   make(Body, 0, len(r)*len(h)),
	}

	return t.Append(r...)
}

// FromCSV returns a new table with data read from a csv reader.
func FromCSV(r *csv.Reader) (*Table, error) {
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

			if ft, err := ParseFTime(lines[i][j]); err == nil {
				r = append(r, ft)
				continue
			}

			r = append(r, lines[i][j])
		}

		t.Append(r)
	}

	return t, nil
}

// FromJSON returns a new table with data parsed from a json-encoded string.
// This string should adhere to the following format.
// 	{"header":["", ...],"types":[0, ...],"body":["", ...]}
func FromJSON(s string) *Table {
	var (
		headerResults = gjson.Get(s, "header").Array()
		n             = len(headerResults)
		h             = make(Header, 0, n)
	)

	for i := 0; i < n; i++ {
		h = append(h, headerResults[i].String())
	}

	t := New(h)
	if 0 < n {
		var (
			typeResults = gjson.Get(s, "types").Array()
			bodyResults = gjson.Get(s, "body").Array()
			mn          = len(bodyResults)
		)

		for i := 0; i < mn; i += n {
			r := make(Row, 0, n)
			for j := 0; j < n; j++ {
				switch Type(typeResults[j].Int()) {
				case Int:
					r = append(r, int(bodyResults[i+j].Int()))
				case Flt:
					r = append(r, float64(bodyResults[i+j].Float()))
				case Bool:
					r = append(r, bool(bodyResults[i+j].Bool()))
				case Time:
					ft, err := ParseFTime(bodyResults[i+j].String())
					if err != nil {
						panic(err.Error())
					}

					r = append(r, ft)
				case Str:
					r = append(r, bodyResults[i+j].String())
				default:
					panic(errType.Error())
				}
			}

			t.Append(r)
		}
	}

	return t
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// Append several rows to a table.
func (t *Table) Append(r ...Row) *Table {
	if 0 < len(r) {
		var i int
		if len(t.body) == 0 {
			if len(t.header) != len(r[0]) {
				panic(errDims.Error())
			}

			for j := 0; j < len(r[0]); j++ {
				switch tp := ParseType(r[0][j]); tp {
				case Int, Flt, Bool, Time, Str:
					t.types = append(t.types, tp)
					t.body = append(t.body, r[0][j])
				default:
					panic(errType.Error())
				}
			}

			i++
		}

		for ; i < len(r); i++ {
			if len(t.header) != len(r[i]) {
				panic(errDims.Error())
			}

			for j := 0; j < len(r[i]); j++ {
				if t.types[j] != ParseType(r[i][j]) {
					panic(errType.Error())
				}

				t.body = append(t.body, r[i][j])
			}
		}
	}

	return t
}

// AppendCol appends a column to a table.
func (t *Table) AppendCol(colName string, c Column) *Table {
	m, n := t.Dims()
	if m != len(c) {
		panic(errDims.Error())
	}

	t.header = append(t.header, colName)
	if 0 < m {
		t.types = append(t.types, ParseType(c[0]))
		if m+len(t.body) <= cap(t.body) {
			t.body = append(t.body, make([]interface{}, m)...)
			for i := m - 1; 0 < i; i-- {
				t.body[i*n+i+n] = c[i]
				copy(t.body[i*n+i:i*n+i+n], t.body[i*n:i*n+n])
			}

			t.body[n] = c[0]
		} else {
			b := make(Body, 0, (m+1)*n)
			for i, mn := 0, m*n; i < mn; i += n {
				b = append(append(b, t.body[i:i+n]...), c[i])
			}

			t.body = b
		}
	}

	return t
}

// Col returns the jth Column.
func (t *Table) Col(j int) Column {
	m, n := t.Dims()
	if n <= j {
		panic(errRange.Error())
	}

	c := make(Column, 0, m)
	for ; j < len(t.body); j += n {
		c = append(c, t.body[j])
	}

	return c
}

// ColBools returns the jth column with each value cast as a boolean.
func (t *Table) ColBools(j int) []bool {
	m, n := t.Dims()
	if n <= j {
		panic(errRange.Error())
	}

	c := make([]bool, 0, m)
	for ; j < len(t.body); j += n {
		c = append(c, t.body[j].(bool))
	}

	return c
}

// ColType returns the type of the jth column.
func (t *Table) ColType(j int) Type {
	return t.types[j]
}

// ColFlts returns the jth column with each value cast as a float.
func (t *Table) ColFlts(j int) []float64 {
	m, n := t.Dims()
	if n <= j {
		panic(errRange.Error())
	}

	c := make([]float64, 0, m)
	for ; j < len(t.body); j += n {
		c = append(c, t.body[j].(float64))
	}

	return c
}

// ColInts returns the jth column with each value cast as an integer.
func (t *Table) ColInts(j int) []int {
	m, n := t.Dims()
	if n <= j {
		panic(errRange.Error())
	}

	c := make([]int, 0, m)
	for ; j < len(t.body); j += n {
		c = append(c, t.body[j].(int))
	}

	return c
}

// ColStrs returns the jth column with each value cast as a string.
func (t *Table) ColStrs(j int) []string {
	m, n := t.Dims()
	if n <= j {
		panic(errRange.Error())
	}

	c := make([]string, 0, m)
	for ; j < len(t.body); j += n {
		c = append(c, t.body[j].(string))
	}

	return c
}

// ColTimes returns the jth column with each value cast as a time object.
func (t *Table) ColTimes(j int) []time.Time {
	m, n := t.Dims()
	if n <= j {
		panic(errRange.Error())
	}

	c := make([]time.Time, 0, m)
	for ; j < len(t.body); j += n {
		c = append(c, t.body[j].(FTime).time)
	}

	return c
}

// ColTypes returns the column types.
func (t *Table) ColTypes() Types {
	return append(make(Types, 0, len(t.types)), t.types...)
}

// Copy a table.
func (t *Table) Copy() *Table {
	cpy := Table{
		header: append(make(Header, 0, len(t.header)), t.header...),
		types:  append(make(Types, 0, len(t.types)), t.types...),
		body:   append(make(Body, 0, len(t.body)), t.body...),
	}

	return &cpy
}

// Dims returns the number of rows and the number of columns in a table body.
func (t *Table) Dims() (int, int) {
	m, n := len(t.body), len(t.header)
	if 0 < n {
		m /= n
	}

	return m, n
}

// Equal determines if two tables are equal.
func (t *Table) Equal(tbl *Table) bool {
	return t.header.Equal(tbl.header) && t.types.Equal(tbl.types) && t.body.Equal(tbl.body)
}

// Get the (i,j)th value.
func (t *Table) Get(i, j int) interface{} {
	return t.body[i*len(t.header)+j]
}

// GetBool the (i,j)th value as a boolean.
func (t *Table) GetBool(i, j int) bool {
	return t.body[i*len(t.header)+j].(bool)
}

// GetFlt the (i,j)th value as a float.
func (t *Table) GetFlt(i, j int) float64 {
	return t.body[i*len(t.header)+j].(float64)
}

// GetInt the (i,j)th value as an integer.
func (t *Table) GetInt(i, j int) int {
	return t.body[i*len(t.header)+j].(int)
}

// GetStr the (i,j)th value as a string.
func (t *Table) GetStr(i, j int) string {
	return t.body[i*len(t.header)+j].(string)
}

// GetTime the (i,j)th value as a time object.
func (t *Table) GetTime(i, j int) time.Time {
	return t.body[i*len(t.header)+j].(FTime).time
}

// Header returns the header.
func (t *Table) Header() Header {
	return append(make(Header, 0, len(t.header)), t.header...)
}

// Insert a row into the ith position.
func (t *Table) Insert(i int, r Row) *Table {
	m, _ := t.Dims()
	return t.Append(r).Swap(i, m)
}

// InsertCol inserts a column into the jth position.
func (t *Table) InsertCol(j int, colName string, c Column) *Table {
	_, n := t.Dims()
	return t.AppendCol(colName, c).SwapCols(j, n)
}

// Join several tables having the same number of rows into one.
func Join(tbl ...*Table) *Table {
	// TODO: It is expensive to repeatedly append columns
	var t *Table
	if 0 < len(tbl) {
		t = tbl[0].Copy()
		m0, _ := t.Dims()
		for i := 1; i < len(tbl); i++ {
			mi, ni := tbl[i].Dims()
			if m0 != mi {
				panic(errDims.Error())
			}

			for j := 0; j < ni; j++ {
				t.AppendCol(tbl[i].header[j], tbl[i].Col(j))
			}
		}
	}

	return t
}

// MarshalJSON ...
func (t *Table) MarshalJSON() ([]byte, error) {
	return []byte(t.ToJSON()), nil
}

// Remove ...
func (t *Table) Remove(i int) Row {
	var (
		_, n = t.Dims()
		r    = NewRow(t.body[i*n : (i+1)*n]...)
	)

	t.body = append(t.body[:i*n], t.body[(i+1)*n:]...)
	if len(t.body) < n {
		t.types = t.types[:0]
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
	t.types = append(t.types[:j], t.types[j+1:]...)
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

	for i, mn := 0, m*n; i < mn; i += n {
		rs = append(rs, NewRow(t.body[i:i+n]...))
	}

	return rs
}

// Set ...
func (t *Table) Set(i, j int, v interface{}) *Table {
	if t.types[j] != ParseType(v) {
		panic(errType.Error())
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
	sb.WriteString("\n" + h + "\n|")
	for j := 0; j < n; j++ {
		switch t.types[j] {
		case Flt, Int:
			sb.WriteString(strings.Repeat(" ", ws[j]-len(t.header[j])+1) + t.header[j] + " |")
		case Bool, Time, Str:
			sb.WriteString(" " + t.header[j] + strings.Repeat(" ", ws[j]-len(t.header[j])+1) + "|")
		default:
		}
	}

	sb.WriteString("\n" + h)
	for i := 0; i < len(b); i += n {
		sb.WriteString("\n|")
		for j, k := 0, i; j < n && k < len(b); j, k = j+1, k+1 {
			switch t.types[j] {
			case Flt, Int:
				// Right-aligned
				sb.WriteString(strings.Repeat(" ", ws[j]-len(b[k])+1) + b[k] + " |")
			case Bool, Time, Str:
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
	for i, mn := 0, m*n; i < mn; i += n {
		r := make([]string, 0, n)
		for j := 0; j < n; j++ {
			switch t.types[j] {
			case Bool:
				r = append(r, strconv.FormatBool(t.body[i+j].(bool)))
			case Flt:
				if v := t.body[i+j].(float64); v == float64(int64(v)) {
					r = append(r, strconv.FormatFloat(v, 'f', 1, 64)) // Forces f.0 when value is an integer
				} else {
					r = append(r, strconv.FormatFloat(v, 'f', -1, 64))
				}
			case Int:
				r = append(r, strconv.Itoa(t.body[i+j].(int)))
			case Time:
				r = append(r, t.body[i+j].(FTime).String())
			case Str:
				r = append(r, t.body[i+j].(string))
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
	t.types[i], t.types[j] = t.types[j], t.types[i]

	_, n := t.Dims()
	for kn := 0; kn < len(t.body); kn += n {
		t.body[kn+i], t.body[kn+j] = t.body[kn+j], t.body[kn+i]
	}

	return t
}

// ToCSV ...
func (t *Table) ToCSV(w *csv.Writer) error {
	return w.WriteAll(t.Strings())
}

// ToJSON ...
func (t *Table) ToJSON() string {
	m, n := t.Dims()
	if n == 0 {
		return `{"header":[],"types":[],"body":[]}`
	}

	var sb strings.Builder
	sb.Grow((m + 3) * n * 8)
	sb.WriteString(`{"header":["` + strings.Join([]string(t.header), `","`) + `"],"types":[` + strconv.Itoa(int(t.types[0])))
	for j := 1; j < n; j++ {
		sb.WriteString(`,` + strconv.Itoa(int(t.types[j])))
	}

	sb.WriteString(`],"body":[`)
	if 0 < m {
		switch t.types[0] {
		case Int:
			sb.WriteString(strconv.FormatInt(int64(t.body[0].(int)), 10))
		case Flt:
			sb.WriteString(strconv.FormatFloat(t.body[0].(float64), 'f', -1, 64))
		case Bool:
			sb.WriteString(strconv.FormatBool(t.body[0].(bool)))
		case Time:
			sb.WriteString(`"` + t.body[0].(FTime).String() + `"`)
		case Str:
			sb.WriteString(`"` + t.body[0].(string) + `"`)
		default:
		}

		for j := 1; j < n; j++ {
			switch t.types[j] {
			case Int:
				sb.WriteString(`,` + strconv.FormatInt(int64(t.body[j].(int)), 10))
			case Flt:
				sb.WriteString(`,` + strconv.FormatFloat(t.body[j].(float64), 'f', -1, 64))
			case Bool:
				sb.WriteString(`,` + strconv.FormatBool(t.body[j].(bool)))
			case Time:
				sb.WriteString(`,"` + t.body[j].(FTime).String() + `"`)
			case Str:
				sb.WriteString(`,"` + t.body[j].(string) + `"`)
			default:
			}
		}

		for i, mn := n, m*n; i < mn; i += n {
			for j := 0; j < n; j++ {
				switch t.types[j] {
				case Int:
					sb.WriteString(`,` + strconv.FormatInt(int64(t.body[i+j].(int)), 10))
				case Flt:
					sb.WriteString(`,` + strconv.FormatFloat(t.body[i+j].(float64), 'f', -1, 64))
				case Bool:
					sb.WriteString(`,` + strconv.FormatBool(t.body[i+j].(bool)))
				case Time:
					sb.WriteString(`,"` + t.body[i+j].(FTime).String() + `"`)
				case Str:
					sb.WriteString(`,"` + t.body[i+j].(string) + `"`)
				default:
				}
			}
		}
	}

	sb.WriteString(`]}`)
	return sb.String()
}

// UnmarshalJSON ...
func (t *Table) UnmarshalJSON(b []byte) error {
	if t1 := FromJSON(string(b)); t1 != nil {
		*t = *t1
	}

	return nil
}
