package table

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

// A Table holds rows/columns of data.
type Table struct {
	Name       string
	Header     Header
	Body       []Row
	Rows, Cols int
}

// A Header describes column data.
type Header []string

// A Row is a single entry in a table.
type Row []interface{}

// A Col is a collection of the ith values in a body of rows.
type Col []interface{}

// New returns an empty table.
func New(tableName string, maxRows, maxCols int) *Table {
	return &Table{
		Name:   tableName,
		Header: make(Header, 0, maxCols),
		Body:   make([]Row, 0, maxRows),
	}
}

// Dims returns the number of rows and columns of a table.
func (t *Table) Dims() (int, int) {
	return t.Rows, t.Cols
}

// AppendRow to a table.
func (t *Table) AppendRow(r Row) {
	t.Body = append(t.Body, r)
	t.Rows++
}

// RemoveRow from a table.
func (t *Table) RemoveRow(i int) Row {
	r := t.Body[i]
	if i+1 == t.Rows {
		t.Body = t.Body[:i]
	} else {
		t.Body = append(t.Body[:i], t.Body[i+1:]...)
	}

	t.Rows--
	return r
}

// AppendCol to a table.
func (t *Table) AppendCol(colName string, c Col) {
	t.Header = append(t.Header, colName)
	n := len(c)
	for t.Rows < n {
		t.AppendRow(make(Row, t.Cols))
	}

	for n < t.Rows {
		c = append(c, nil)
		n++
	}

	for i := range t.Body {
		t.Body[i] = append(t.Body[i], c[i])
	}

	t.Cols++
}

// RemoveCol from a table.
func (t *Table) RemoveCol(i int) (string, Col) {
	h := t.Header[i]
	c := make(Col, 0, t.Rows)
	if i+1 == t.Cols {
		t.Header = t.Header[:i]
		for j := range t.Body {
			c = append(c, t.Body[j][i])
			t.Body[j] = t.Body[j][:i]
		}
	} else {
		t.Header = append(t.Header[:i], t.Header[i+1:]...)
		for j := range t.Body {
			c = append(c, t.Body[j][i])
			t.Body[j] = append(t.Body[j][:i], t.Body[j][:i+1]...)
		}
	}

	t.Cols--
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
