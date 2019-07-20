package table

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// A Table holds rows/columns of data.
type Table struct {
	Name          string
	header        Header
	body          []Row
	rows, columns int
}

// A Header describes column data.
type Header []string

// A Row is a single entry in a table.
type Row []interface{}

// A Column is a collection of the ith values in a body of rows.
type Column []interface{}

// New returns an empty table.
func New(tableName string, maxRows, maxColumns int) *Table {
	return &Table{
		Name:   tableName,
		header: make(Header, 0, maxColumns),
		body:   make([]Row, 0, maxRows),
	}
}

// Get returns the (i,j)th value.
func (t *Table) Get(i, j int) interface{} {
	return t.body[i][j]
}

// Set the (i,j)th cell to a given value.
func (t *Table) Set(v interface{}, i, j int) {
	t.body[i][j] = v
}

// GetColumnHeader at a given index.
func (t *Table) GetColumnHeader(i int) string {
	return t.header[i]
}

// SetColumnHeader to a given value.
func (t *Table) SetColumnHeader(columnHeader string, i int) {
	t.header[i] = columnHeader
}

// GetRow returns a copy of a row.
func (t *Table) GetRow(i int) Row {
	return t.body[i].Copy()
}

// GetColumn returns the column header and a copy of the column at a given index.
func (t *Table) GetColumn(i int) (string, Column) {
	c := make(Column, 0, len(t.body))
	for j := range t.body {
		c = append(c, t.body[j][i])
	}

	return t.header[i], c
}

// Dimensions returns the number of rows and columns of a table.
func (t *Table) Dimensions() (int, int) {
	return t.rows, t.columns
}

// AppendRow to a table.
func (t *Table) AppendRow(r Row) {
	t.body = append(t.body, r)
	t.rows++

	n := len(r)
	if t.columns < n {
		t.columns = n
	}

	for len(t.header) < n {
		t.header = append(t.header, "")
	}

	for i := range t.body {
		for len(t.body[i]) < n {
			t.body[i] = append(t.body[i], nil)
		}
	}
}

// RemoveRow from a table.
func (t *Table) RemoveRow(i int) Row {
	r := t.body[i]
	if i+1 == t.rows {
		t.body = t.body[:i]
	} else {
		t.body = append(t.body[:i], t.body[i+1:]...)
	}

	t.rows--
	return r
}

// AppendColumn to a table.
func (t *Table) AppendColumn(columnHeader string, c Column) {
	t.header = append(t.header, columnHeader)
	n := len(c)
	for t.rows < n {
		t.AppendRow(make(Row, t.columns))
	}

	for n < t.rows {
		c = append(c, nil)
		n++
	}

	for i := range t.body {
		t.body[i] = append(t.body[i], c[i])
	}

	t.columns++
}

// RemoveColumn from a table.
func (t *Table) RemoveColumn(i int) (string, Column) {
	h := t.header[i]
	c := make(Column, 0, t.rows)
	if i+1 == t.columns {
		t.header = t.header[:i]
		for j := range t.body {
			c = append(c, t.body[j][i])
			t.body[j] = t.body[j][:i]
		}
	} else {
		t.header = append(t.header[:i], t.header[i+1:]...)
		for j := range t.body {
			c = append(c, t.body[j][i])
			t.body[j] = append(t.body[j][:i], t.body[j][:i+1]...)
		}
	}

	t.columns--

	// Remove empty rows
	for j := 0; j < t.rows; j++ {
		if t.body[j].isEmpty() {
			if j+1 == t.rows {
				t.body = t.body[:j]
			} else {
				t.body = append(t.body[:j], t.body[j+1:]...)
			}

			t.rows--
		}
	}

	return h, c
}

// ImportCSV imports a csv file into a table and returns it.
func ImportCSV(path, tableName string) (*Table, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var (
		reader = csv.NewReader(file)
		t      = New(tableName, 0, 0)
		value  float64
		line   []string
	)
	line, err = reader.Read()
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
		return t, nil
	}

	t.header = Header(line)
	for {
		line, err = reader.Read()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return t, nil
		}

		r := make(Row, 0, len(line))
		for _, text := range line {
			value, err = strconv.ParseFloat(text, 64)
			if err != nil {
				r = append(r, text)
			} else {
				r = append(r, value)
			}
		}

		t.AppendRow(r)
	}
}

// Copy a row.
func (r Row) Copy() Row {
	s := make(Row, len(r))
	copy(s, r)
	return s
}

// Copy a column.
func (c Column) Copy() Column {
	d := make(Column, len(c))
	copy(d, c)
	return d
}

// isEmpty determines if a row contains data.
func (r Row) isEmpty() bool {
	for i := range r {
		if r[i] != nil {
			return false
		}
	}

	return true
}

func (r Row) String() string {
	sb := strings.Builder{}
	for i := range r {
		if v, ok := r[i].(float64); ok {
			sb.WriteString(fmt.Sprintf("%0.2f", v))
		} else {
			sb.WriteString(fmt.Sprintf("%v", r[i]))
		}
	}

	return sb.String()
}

func (h Header) String() string {
	return strings.Join(h, " ")
}

func (t *Table) FmtString(precision int) string {
	sb := strings.Builder{}
	// widths := make([]int, t.columns)
	// for i := range t.body {
	// 	for j := range t.body[i] {
	// 		v,ok
	// 	}
	// }
	return sb.String()
}
