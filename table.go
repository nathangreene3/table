package table

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type format byte

const (
	unset format = iota
	strFmt
	fltFmt
	intFmt
)

// A Table holds rows/columns of data.
type Table struct {
	Name          string
	header        Header
	body          []Row
	rows, columns int
	formats       []format
	widths        []int
	fltFmt        byte
	fltPrec       int
}

// A Header describes column data.
type Header []string

// A Row is a single entry in a table.
type Row []interface{}

// A Column is a collection of the ith values in a body of rows.
type Column []interface{}

// New returns an empty table.
func New(tableName string, fltFmt byte, fltPrec, maxRows, maxColumns int) *Table {
	return &Table{
		Name:    tableName,
		header:  make(Header, 0, maxColumns),
		body:    make([]Row, 0, maxRows),
		formats: make([]format, 0, maxColumns),
		widths:  make([]int, 0, maxColumns),
		fltFmt:  fltFmt,
		fltPrec: fltPrec,
	}
}

// SetHeader ...
func (t *Table) SetHeader(h Header) {
	n := len(h)
	t.setColumns(n)

	var width int
	for i := 0; i < n; i++ {
		t.header[i] = strings.TrimSpace(h[i])
		width = len(h[i])
		if t.widths[i] < width {
			t.widths[i] = width
		}
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
	var (
		ok       bool
		width    int
		intValue int
		fltValue float64
		n        = len(r)
	)
	t.setColumns(n)
	t.body = append(t.body, r)
	t.rows++

	for i := 0; i < n; i++ {
		switch t.formats[i] {
		case intFmt:
			if intValue, ok = r[i].(int); ok {
				width = len(strconv.Itoa(intValue))
			} else if _, ok = r[i].(float64); ok {
				t.formats[i] = fltFmt
				width = len(strconv.FormatFloat(fltValue, t.fltFmt, t.fltPrec, 64))
			} else {
				t.formats[i] = strFmt
				width = len(r[i].(string))
			}
		case fltFmt:
			if _, ok = r[i].(float64); ok {
				width = len(strconv.FormatFloat(fltValue, t.fltFmt, t.fltPrec, 64))
			} else {
				t.formats[i] = strFmt
				width = len(r[i].(string))
			}
		case strFmt:
			width = len(r[i].(string))
		}

		if t.widths[i] < width {
			t.widths[i] = width
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

// setColumns ...
func (t *Table) setColumns(n int) {
	t.columns = n
	for len(t.header) < n {
		t.header = append(t.header, "")
	}

	for i := range t.body {
		for len(t.body[i]) < n {
			t.body[i] = append(t.body[i], nil)
		}
	}

	for len(t.formats) < n {
		t.formats = append(t.formats, unset)
	}

	for len(t.widths) < n {
		t.widths = append(t.widths, 0)
	}
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
func ImportCSV(path, tableName string, fltFmt byte, fltPrec int) (*Table, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var (
		reader = csv.NewReader(file)
		t      = New(tableName, fltFmt, fltPrec, 0, 0)
		line   []string
	)
	line, err = reader.Read()
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
		return t, nil
	}

	t.SetHeader(line)
	for {
		line, err = reader.Read()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return t, nil
		}

		r := make(Row, 0, len(line))
		for _, strValue := range line {
			r = append(r, strings.TrimSpace(strValue))
		}

		t.AppendRow(r)
	}
}

// ExportCSV ... TODO
func (t *Table) ExportCSV(path string) error {
	return nil
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

// isEmpty determines if a row contains data or not.
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

// String ...
func (t *Table) String() string {
	// Create horizontal line
	sb := strings.Builder{}
	for i := range t.widths {
		sb.WriteString("+" + strings.Repeat("-", t.widths[i]))
	}

	sb.WriteString("+\n")
	hLine := sb.String()
	sb.Reset()

	// Write header
	sb.WriteString(hLine)
	for i := 0; i < t.columns; i++ {
		switch t.formats[i] {
		case intFmt:
			fallthrough
		case fltFmt:
			sb.WriteString("|" + t.header[i] + strings.Repeat(" ", t.widths[i]-len(t.header[i])))
		case strFmt:
			sb.WriteString("|" + strings.Repeat(" ", t.widths[i]-len(t.header[i])) + t.header[i])
		}
	}

	sb.WriteString("|\n")
	sb.WriteString(hLine)

	// Write body
	var strValue string
	for i := 0; i < t.rows; i++ {
		for j := 0; j < t.columns; j++ {
			switch t.formats[i] {
			case intFmt:
				strValue = strconv.Itoa(t.body[i][j].(int))
				sb.WriteString("|" + strValue + strings.Repeat(" ", t.widths[i]-len(strValue)))
			case fltFmt:
				strValue = strconv.FormatFloat(t.body[i][j].(float64), t.fltFmt, t.fltPrec, 64)
				sb.WriteString("|" + strValue + strings.Repeat(" ", t.widths[i]-len(strValue)))
			case strFmt:
				strValue = t.body[i][j].(string)
				sb.WriteString("|" + strValue + strings.Repeat(" ", t.widths[i]-len(strValue)))
			}
		}

		sb.WriteString("|\n")
	}

	sb.WriteString(hLine)
	return sb.String()
}
