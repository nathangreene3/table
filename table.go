package table

import (
	"bytes"
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type (
	// A Table holds tabular data.
	Table struct {
		header Header
		types  Types
		body   Body
	}

	// Filterer determines the criteria for retaining a row.
	Filterer func(r Row) bool

	// Generator defines the (i,j)th value.
	Generator func(i, j int) interface{}

	// Mapper mutates a given row.
	Mapper func(r Row)

	// Reducer combines two rows into one. Only the destination row
	// should be mutated.
	Reducer func(dst, src Row)
)

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

// FromCSV returns a new table with data read from a csv file.
func FromCSV(fileName string) (*Table, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	lines, err := csv.NewReader(bytes.NewReader(b)).ReadAll()
	if err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return New(nil), nil
	}

	t := Table{
		header: Header(lines[0]),
		types:  make(Types, 0, len(lines[0])),
		body:   make(Body, 0, len(lines[0])*len(lines[1:])),
	}

	for i := 1; i < len(lines); i++ {
		r := make(Row, 0, len(lines[i]))
		for j := 0; j < len(lines[i]); j++ {
			if n, err := strconv.ParseInt(lines[i][j], 10, strconv.IntSize); err == nil {
				r = append(r, int(n))
			} else if f, err := strconv.ParseFloat(lines[i][j], strconv.IntSize); err == nil {
				r = append(r, f)
			} else if b, err := strconv.ParseBool(lines[i][j]); err == nil {
				r = append(r, b)
			} else if ft, err := ParseFTime(lines[i][j]); err == nil {
				r = append(r, ft)
			} else {
				r = append(r, lines[i][j])
			}
		}

		t.Append(r)
	}

	return &t, nil
}

// FromJSON returns a new table with data parsed from a json-encoded
// string. This string should adhere to the following format.
// 	{"header":["", ...],"types":[0, ...],"body":["", ...]}
func FromJSON(s string) (*Table, error) {
	var (
		header = gjson.Get(s, "header").Array()
		types  = gjson.Get(s, "types").Array()
		body   = gjson.Get(s, "body").Array()
		n      = len(header)
		mn     = len(body)
		t      = Table{
			header: make(Header, 0, n),
			types:  make(Types, 0, n),
			body:   make(Body, 0, mn),
		}
	)

	for i := 0; i < n; i++ {
		t.header = append(t.header, header[i].String())
	}

	for i := 0; i < mn; i += n {
		r := make(Row, 0, n)
		for j := 0; j < n; j++ {
			switch Type(types[j].Int()) {
			case Int:
				r = append(r, int(body[i+j].Int()))
			case Flt:
				r = append(r, float64(body[i+j].Float()))
			case Bool:
				r = append(r, bool(body[i+j].Bool()))
			case Time:
				ft, err := ParseFTime(body[i+j].String())
				if err != nil {
					return nil, err
				}

				r = append(r, ft)
			case Str:
				r = append(r, body[i+j].String())
			default:
				return nil, errors.New(errType)
			}
		}

		t.Append(r)
	}

	return &t, nil
}

// Generate returns a new table generated by a generator.
func Generate(h Header, m int, f Generator) *Table {
	t := Table{
		header: h.Copy(),
		types:  make(Types, 0, len(h)),
		body:   make(Body, 0, m*len(h)),
	}

	if 0 < m {
		for j := 0; j < len(h); j++ {
			t.body = append(t.body, f(0, j))
			t.types = append(t.types, Parse(t.body[j]))
		}

		for i := 1; i < m; i++ {
			for j := 0; j < len(h); j++ {
				t.body = append(t.body, f(i, j))
				if Parse(t.body[i*len(h)+j]) != t.types[j] {
					panic(errType)
				}
			}
		}
	}

	return &t
}

// --------------------------------------------------------------------
// Methods
// --------------------------------------------------------------------

// Append several rows to a table.
func (t *Table) Append(r ...Row) *Table {
	if 0 < len(r) {
		var i int
		if len(t.body) == 0 {
			if len(t.header) != len(r[0]) {
				panic(errDims)
			}

			for j := 0; j < len(r[0]); j++ {
				switch tp := Parse(r[0][j]); tp {
				case Int, Flt, Bool, Time, Str:
					t.types = append(t.types, tp)
					t.body = append(t.body, r[0][j])
				default:
					panic(errType)
				}
			}

			i++
		}

		for ; i < len(r); i++ {
			if len(t.header) != len(r[i]) {
				panic(errDims)
			}

			for j := 0; j < len(r[i]); j++ {
				if t.types[j] != Parse(r[i][j]) {
					panic(errType)
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
	if n == 0 && 0 < len(c) {
		tp := c.Type()
		if tp == 0 {
			panic(errType)
		}

		t.header = Header{colName}
		t.types = []Type{tp}
		t.body = append(make(Body, 0, len(c)), c...)
		return t
	}

	if m != len(c) {
		panic(errDims)
	}

	t.header = append(t.header, colName)
	if 0 < m {
		t.types = append(t.types, Parse(c[0]))
		if m+len(t.body) <= cap(t.body) {
			t.body = append(t.body, make([]interface{}, m)...)
			for i := m - 1; 0 < i; i-- {
				j0 := i * n
				j1 := j0 + i
				j2 := j1 + n
				t.body[j2] = c[i]
				copy(t.body[j1:j2], t.body[j0:j0+n])
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

// Bool returns the (i,j)th value as a boolean.
func (t *Table) Bool(i, j int) bool {
	return t.body[i*len(t.header)+j].(bool)
}

// Col returns the jth Column.
func (t *Table) Col(j int) Column {
	m, n := t.Dims()
	if n <= j {
		panic(errRange)
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
		panic(errRange)
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

// ColFloats returns the jth column with each value cast as a float.
func (t *Table) ColFloats(j int) []float64 {
	m, n := t.Dims()
	if n <= j {
		panic(errRange)
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
		panic(errRange)
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
		panic(errRange)
	}

	c := make([]string, 0, m)
	for ; j < len(t.body); j += n {
		c = append(c, t.body[j].(string))
	}

	return c
}

// ColTimes returns the jth column with each value cast as a time
// object.
func (t *Table) ColTimes(j int) []time.Time {
	m, n := t.Dims()
	if n <= j {
		panic(errRange)
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

// Dims returns the number of rows and the number of columns in a
// table body.
func (t *Table) Dims() (int, int) {
	m, n := len(t.body), len(t.header)
	if 0 < n {
		m /= n
	}

	return m, n
}

// Equal determines if two tables are equal.
func (t *Table) Equal(tbl *Table) bool {
	return t == tbl || (t.header.Equal(tbl.header) && t.types.Equal(tbl.types) && t.body.Equal(tbl.body))
}

// Filter applies a filterer on each row. Each row in which f
// evaluates as true is retained; all other rows are discarded.
func (t *Table) Filter(f Filterer) *Table {
	m, n := t.Dims()
	for i := (m - 1) * n; 0 <= i; i -= n {
		if !f(Row(t.body[i : i+n])) {
			t.body = append(t.body[:i], t.body[i+n:]...)
		}
	}

	if len(t.body) < n {
		t.types = t.types[:0]
	}

	return t
}

// Float returns the (i,j)th value as a float.
func (t *Table) Float(i, j int) float64 {
	return t.body[i*len(t.header)+j].(float64)
}

// Format returns a formatted table given format rules.
func (t *Table) Format(fmt Format) string {
	var (
		m, n = t.Dims()
		mn   = m * n
		ws   = make([]int, 0, n) // Column widths
	)

	for j := 0; j < n; j++ {
		ws = append(ws, len(t.header[j]))
	}

	b := t.body.Strings()
	for i := 0; i < mn; i += n {
		for j := 0; j < n; j++ {
			if w := len(b[i+j]); ws[j] < w {
				ws[j] = w
			}
		}
	}

	var sb strings.Builder
	sb.Grow(256) // TODO: Estimate how big a table may be
	sb.WriteByte('\n')
	if 0 < len(fmt.UpperHoriz) {
		sb.WriteString(fmt.UpperLeftHorizDelim)
		if 0 < n {
			sb.WriteString(strings.Repeat(fmt.UpperHoriz, ws[0]+2))
		}

		for j := 1; j < n; j++ {
			sb.WriteString(fmt.UpperMidHorizDelim + strings.Repeat(fmt.UpperHoriz, ws[j]+2))
		}

		sb.WriteString(fmt.UpperRightHorizDelim + "\n")
	}

	if 0 < len(fmt.HeaderLeftDelim) {
		sb.WriteString(fmt.HeaderLeftDelim)
	}

	if 0 < n {
		switch t.types[0] {
		case Flt, Int:
			sb.WriteString(strings.Repeat(" ", ws[0]-len(t.header[0])+1) + t.header[0] + " ")
		case Bool, Time, Str:
			sb.WriteString(" " + t.header[0] + strings.Repeat(" ", ws[0]-len(t.header[0])+1))
		default:
			panic(errType)
		}
	}

	for j := 1; j < n; j++ {
		switch t.types[j] {
		case Flt, Int:
			sb.WriteString(fmt.HeaderMidDelim + strings.Repeat(" ", ws[j]-len(t.header[j])+1) + t.header[j] + " ")
		case Bool, Time, Str:
			sb.WriteString(fmt.HeaderMidDelim + " " + t.header[j] + strings.Repeat(" ", ws[j]-len(t.header[j])+1))
		default:
			panic(errType)
		}
	}

	sb.WriteString(fmt.HeaderRightDelim + "\n")
	if 0 < len(fmt.MiddleHoriz) {
		sb.WriteString(fmt.MiddleLeftHorizDelim)
		if 0 < n {
			sb.WriteString(strings.Repeat(fmt.MiddleHoriz, ws[0]+2))
		}

		for j := 1; j < n; j++ {
			sb.WriteString(fmt.MiddleMidHorizDelim + strings.Repeat(fmt.MiddleHoriz, ws[j]+2))
		}

		sb.WriteString(fmt.MiddleRightHorizDelim + "\n")
	}

	for i := 0; i < mn; i += n {
		switch t.types[0] {
		case Flt, Int:
			sb.WriteString(fmt.RowLeftDelim + strings.Repeat(" ", ws[0]-len(b[i])+1) + b[i] + " ")
		case Bool, Time, Str:
			sb.WriteString(fmt.RowLeftDelim + " " + b[i] + strings.Repeat(" ", ws[0]-len(b[i])+1))
		default:
			panic(errType)
		}

		for j := 1; j < n; j++ {
			ij := i + j
			switch t.types[j] {
			case Flt, Int:
				sb.WriteString(fmt.RowMidDelim + strings.Repeat(" ", ws[j]-len(b[ij])+1) + b[ij] + " ")
			case Bool, Time, Str:
				sb.WriteString(fmt.RowMidDelim + " " + b[ij] + strings.Repeat(" ", ws[j]-len(b[ij])+1))
			default:
				panic(errType)
			}
		}

		sb.WriteString(fmt.RowRightDelim + "\n")
	}

	// Bottom horizontal line
	if 0 < len(fmt.BottomHoriz) {
		sb.WriteString(fmt.BottomLeftHorizDelim)
		if 0 < n {
			sb.WriteString(strings.Repeat(fmt.BottomHoriz, ws[0]+2))
		}

		for j := 1; j < n; j++ {
			sb.WriteString(fmt.BottomMidHorizDelim + strings.Repeat(fmt.BottomHoriz, ws[j]+2))
		}

		sb.WriteString(fmt.BottomRightHorizDelim + "\n")
	}

	return sb.String()
}

// Header returns the header.
func (t *Table) Header() Header {
	return append(make(Header, 0, len(t.header)), t.header...)
}

// Insert a row into the ith position.
func (t *Table) Insert(i int, r Row) *Table {
	m, n := t.Dims()
	for i := 0; i < n; i++ {
		t.body = append(t.body, nil)
	}

	var (
		j  = i * n
		k  = j + n
		mn = m * n
	)

	copy(t.body[k:mn+n], t.body[j:mn])
	copy(t.body[j:k], r)
	return t.Append(r).Swap(i, m)
}

// InsertCol inserts a column into the jth position.
func (t *Table) InsertCol(j int, colName string, c Column) *Table {
	return t.AppendCol(colName, c).SwapCols(j, len(t.header))
}

// Int returns the (i,j)th value as an integer.
func (t *Table) Int(i, j int) int {
	return t.body[i*len(t.header)+j].(int)
}

// Join several tables having the same number of rows into one.
func Join(tbl ...*Table) *Table {
	if len(tbl) == 0 {
		return New(NewHeader())
	}

	m0, n := tbl[0].Dims()
	ns := append(make([]int, 0, len(tbl)), n)
	for i := 1; i < len(tbl); i++ {
		mi, ni := tbl[i].Dims()
		if m0 != mi {
			panic(errDims)
		}

		n += ni
		ns = append(ns, ni)
	}

	h := make(Header, 0, n)
	for i := 0; i < len(tbl); i++ {
		h = append(h, tbl[i].header...)
	}

	f := func(i, j int) interface{} {
		// TODO: Computing k and modifying j is slow. Figure out how to map j to k and j to ns[k].
		var k int
		for ; k < len(tbl) && ns[k] <= j; k++ {
			j -= ns[k]
		}

		return tbl[k].body[i*ns[k]+j]
	}

	return Generate(h, m0, f)
}

// JSON returns a json-encoded string representing a table.
func (t *Table) JSON() string {
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

// Map mutates each row in a table and updates the column types.
func (t *Table) Map(f Mapper) *Table {
	n := len(t.header)
	f(Row(t.body[:n]))

	var i int
	for ; i < n; i++ {
		t.types[i] = Parse(t.body[i])
	}

	for ; i < len(t.body); i += n {
		f(Row(t.body[i : i+n]))
	}

	return t
}

// MarshalJSON returns a list of json-encoded bytes. This implements
// the json.Marshaller interface.
func (t *Table) MarshalJSON() ([]byte, error) {
	return []byte(t.JSON()), nil
}

// Reduce returns a row that is the product of applying a reducer on
// each row in a table. A copy of the first row is used as the
// accumulator.
func (t *Table) Reduce(f Reducer) Row {
	var (
		m, n = t.Dims()
		r    = make(Row, 0, n)
	)

	if 0 < m {
		r = append(r, t.body[:n]...)
	}

	for i := n; i < len(t.body); i += n {
		f(r, Row(t.body[i:i+n]))
	}

	return r
}

// Remove removes and returns the ith row from a table.
func (t *Table) Remove(i int) Row {
	var (
		n = len(t.header)
		r = NewRow(t.body[i*n : (i+1)*n]...)
	)

	t.body = append(t.body[:i*n], t.body[(i+1)*n:]...)
	if len(t.body) < n {
		t.types = t.types[:0]
	}

	return r
}

// RemoveCol removes and returns the jth column from a table.
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

// Row returns the ith row from a table.
func (t *Table) Row(i int) Row {
	n := len(t.header)
	return NewRow(t.body[i*n : (i+1)*n]...)
}

// Rows returns a list of all the rows in a table.
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

// Set the (i,j)th value in a table.
func (t *Table) Set(i, j int, v interface{}) *Table {
	if t.types[j] != Parse(v) {
		panic(errType)
	}

	t.body[i*len(t.header)+j] = v
	return t
}

// SetColName sets the jth column name.
func (t *Table) SetColName(j int, s string) *Table {
	t.header[j] = s
	return t
}

// SetHeader sets the header.
func (t *Table) SetHeader(h Header) *Table {
	if len(h) != len(t.header) {
		panic(errDims)
	}

	copy(t.header, h)
	return t
}

// Sort sorts a table on the jth column.
func (t *Table) Sort(j int) *Table {
	return t.Stable(j) // TODO: Replace.
}

// Stable sorts a table on the jth column.
func (t *Table) Stable(j int) *Table {
	m, n := t.Dims()
	for k := 1; k < m; k++ {
		for i := k - 1; 0 <= i; i-- {
			switch t.types[j] {
			case Int:
				if t.body[i*n+j].(int) <= t.body[(i+1)*n+j].(int) {
					i = -1
				}
			case Flt:
				if t.body[i*n+j].(float64) <= t.body[(i+1)*n+j].(float64) {
					i = -1
				}
			case Bool:
				if !t.body[i*n+j].(bool) || t.body[(i+1)*n+j].(bool) {
					i = -1
				}
			case Time:
				if t.body[i*n+j].(FTime).Compare(t.body[(i+1)*n+j].(FTime)) <= 0 {
					i = -1
				}
			case Str:
				if t.body[i*n+j].(string) <= t.body[(i+1)*n+j].(string) {
					i = -1
				}
			default:
				panic(errType)
			}

			if 0 <= i {
				t.Swap(i, i+1)
			}
		}
	}

	return t
}

// Str returns the (i,j)th value as a string.
func (t *Table) Str(i, j int) string {
	return t.body[i*len(t.header)+j].(string)
}

// String returns a string representing a table.
func (t *Table) String() string {
	return "[" + t.header.String() + " | " + t.body.String() + "]"
}

// Strings returns a list of string lists. The first string list is
// the header.
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

// Swap swaps two rows in a table.
func (t *Table) Swap(i, j int) *Table {
	for k, n := 0, len(t.header); k < n; k++ {
		ik, jk := i*n+k, j*n+k
		t.body[ik], t.body[jk] = t.body[jk], t.body[ik]
	}

	return t
}

// SwapCols swaps two columns in a table.
func (t *Table) SwapCols(i, j int) *Table {
	t.header[i], t.header[j] = t.header[j], t.header[i]
	t.types[i], t.types[j] = t.types[j], t.types[i]

	for kn, n := 0, len(t.header); kn < len(t.body); kn += n {
		t.body[kn+i], t.body[kn+j] = t.body[kn+j], t.body[kn+i]
	}

	return t
}

// Time the (i,j)th value as a time object.
func (t *Table) Time(i, j int) time.Time {
	return t.body[i*len(t.header)+j].(FTime).time
}

// UnmarshalJSON reads a list of json-encoded bytes into a table.
// Implements the json.Unmarshaller interface.
func (t *Table) UnmarshalJSON(b []byte) error {
	t1, err := FromJSON(string(b))
	if err == nil {
		*t = *t1
	}

	return err
}

// Validate returns an error if a table is in an invalid state.
func (t *Table) Validate() error {
	m, n := t.Dims()
	if len(t.types) != n || len(t.body) != m*n {
		return errors.New(errDims)
	}

	for i := 0; i < m; i++ {
		for k, j := i*n, 0; j < n; j++ {
			if Parse(t.body[k+j]) != t.types[j] {
				return errors.New(errType)
			}
		}
	}

	return nil
}

// Value returns the (i,j)th value.
func (t *Table) Value(i, j int) interface{} {
	return t.body[i*len(t.header)+j]
}

// WriteCSV writes a table to a csv file.
func (t *Table) WriteCSV(file string) error {
	file = filepath.Clean(file)
	if !strings.EqualFold(filepath.Ext(file), ".csv") {
		file += ".csv"
	}

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()
	return csv.NewWriter(f).WriteAll(t.Strings())
}
