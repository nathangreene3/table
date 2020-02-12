package table

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

// A Table holds rows/columns of data.
type Table struct {
	header         Header
	body           Body
	colBaseTypes   []baseType
	colWidths      []int
	Name           string
	rows, columns  int
	floatPrecision FltPrecFmt
	floatFmt       FltFmt
}

// FltFmt defines updateBaseTypesAndWidthsting float values.
type FltFmt byte

// FltPrecFmt defines the number of decimal positions displayed
// for float values.
type FltPrecFmt int

const (
	// FltFmtBinExp updateBaseTypesAndWidthss floats as a binary exponent value -dddp±ddd.
	FltFmtBinExp FltFmt = 'b'
	// FltFmtDecExp updateBaseTypesAndWidthss floats as a decimal exponent value (scientific notation) -d.ddde±ddd.
	FltFmtDecExp FltFmt = 'e'
	// FltFmtNoExp updateBaseTypesAndWidthss floats as a decimal value -ddd.ddd.
	FltFmtNoExp FltFmt = 'f'
	// FltFmtLrgExp updateBaseTypesAndWidthss floats as a large exponent value (scientific notation) -d.ddde±ddd.
	FltFmtLrgExp FltFmt = 'g'
)

// New returns an empty table.
func New(name string, floatFmt FltFmt, floatPrec FltPrecFmt, rows ...Row) *Table {
	if floatFmt == 0 {
		floatFmt = FltFmtNoExp
	}

	t := &Table{
		Name:           name,
		header:         make(Header, 0),
		body:           make(Body, 0, len(rows)),
		colBaseTypes:   make([]baseType, 0),
		colWidths:      make([]int, 0),
		floatPrecision: floatPrec,
		floatFmt:       floatFmt,
	}

	return t.AppendRows(rows...)
}

// AppendColumn to a table. If the column header and column is empty, nothing
// happens.
func (t *Table) AppendColumn(columnHeader string, c Column) *Table {
	if len(columnHeader) == 0 && c.isEmpty() {
		return t
	}

	// Increase body size to column size
	n := len(c)
	for t.rows < n {
		t.AppendRow(make(Row, t.columns))
	}

	// Increase column size to body size
	for n < t.rows {
		c = append(c, nil)
		n++
	}

	t.header = append(t.header, strings.TrimSpace(columnHeader))
	for i := range t.body {
		t.body[i] = append(t.body[i], c[i])
	}

	t.columns++
	t.colBaseTypes[t.columns-1] = t.minColBaseType(t.columns - 1)
	return t
}

// AppendRow to a table.
func (t *Table) AppendRow(r Row) *Table {
	if r.isEmpty() {
		return t
	}

	n := len(r)
	t.setColSize(n)
	t.body = append(t.body, r)
	t.rows++

	var w int
	for i := 0; i < n; i++ {
		switch t.colBaseTypes[i] {
		case integerType:
			switch baseTypeOf(r[i]) {
			case integerType:
				w = len(strconv.Itoa(r[i].(int)))
			case floatType:
				t.colBaseTypes[i] = floatType
				w = len(strconv.FormatFloat(r[i].(float64), byte(t.floatFmt), int(t.floatPrecision), 64))
			case stringType:
				t.colBaseTypes[i] = stringType
				w = len(r[i].(string))
			default:
				panic("unknown type")
			}
		case floatType:
			switch baseTypeOf(r[i]) {
			case integerType:
				w = len(strconv.FormatFloat(float64(r[i].(int)), byte(t.floatFmt), int(t.floatPrecision), 64))
			case floatType:
				w = len(strconv.FormatFloat(r[i].(float64), byte(t.floatFmt), int(t.floatPrecision), 64))
			case stringType:
				t.colBaseTypes[i] = stringType
				w = len(r[i].(string))
			default:
				panic("unknown type")
			}
		case stringType:
			switch baseTypeOf(r[i]) {
			case integerType:
				w = len(strconv.Itoa(r[i].(int)))
			case floatType:
				w = len(strconv.FormatFloat(r[i].(float64), byte(t.floatFmt), int(t.floatPrecision), 64))
			case stringType:
				t.colBaseTypes[i] = stringType
				w = len(r[i].(string))
			default:
				panic("unknown type")
			}
		}

		if t.colWidths[i] < w {
			t.colWidths[i] = w
		}
	}

	return t
}

// AppendRows ...
func (t *Table) AppendRows(rows ...Row) *Table {
	for i := range rows {
		t.AppendRow(rows[i])
	}

	return t
}

// Clean removes empty rows and columns.
func (t *Table) Clean() *Table {
	// Remove empty rows
	for i := t.rows - 1; 0 <= i; i-- {
		if t.body[i].isEmpty() {
			t.RemoveRow(i)
		}
	}

	// Remove empty columns. Named columns are ignored.
	for j := t.columns - 1; 0 <= j; j-- {
		if isEmpty := len(t.header[j]) == 0; isEmpty {
			for i := 0; i < t.rows && isEmpty; i++ {
				isEmpty = t.body[i][j] == nil
			}

			if isEmpty {
				t.RemoveColumn(j)
			}
		}
	}

	return t.updateBaseTypesAndWidths()
}

// Column returns the column header and a copy of the column at a
// given index.
func (t *Table) Column(i int) (string, Column) {
	c := make(Column, 0, len(t.body))
	for _, r := range t.body {
		c = append(c, r[i])
	}

	return t.header[i], c
}

// ColumnHeader at a given index.
func (t *Table) ColumnHeader(i int) string {
	return t.header[i]
}

// TODO: t.Compare

// Copy a table.
func (t *Table) Copy() *Table {
	cpy := New(t.Name, t.floatFmt, t.floatPrecision, t.body...).SetHeader(t.header)
	return cpy
}

// Dimensions returns the number of rows and columns of a table.
func (t *Table) Dimensions() (int, int) {
	return t.rows, t.columns
}

// updateBaseTypesAndWidths a table. This updates each column base type to its
// weakest base type and updates each column width to the largest
// each needs to be when updateBaseTypesAndWidthsed as a string.
func (t *Table) updateBaseTypesAndWidths() *Table {
	// Reset each column updateBaseTypesAndWidths as an integer and to be as wide as the column
	// header. Each column base type is used to determine alignment and doesn't
	// affect the updateBaseTypesAndWidthsting of each (i,j)th value.
	for i := 0; i < t.columns; i++ {
		t.colBaseTypes[i] = integerType
		t.colWidths[i] = len(t.header[i])
	}

	var (
		bt baseType
		w  int
	)

	for _, r := range t.body {
		for j, v := range r {
			// Set the jth column base type to the minimum (weakest) base type
			// found in each value at (i,j). Recall from base_type.go that
			// strings < floats < ints. Then update the jth column width
			// depending on the weakest base type found.

			// minColBaseType is not appropriate here because the (i,j)th value
			// is not reset in dependence on the minimum column base type
			// returned.
			if bt = baseTypeOf(v); bt < t.colBaseTypes[j] {
				t.colBaseTypes[j] = bt
			}

			// Determine the width of the (i,j)th value when converted to its base type and update the column width if the column width is too small to support it.
			switch bt {
			case integerType:
				w = len(strconv.Itoa(v.(int)))
			case floatType:
				w = len(strconv.FormatFloat(v.(float64), byte(t.floatFmt), int(t.floatPrecision), 64))
			case stringType:
				w = len(v.(string))
			default:
				panic("unknown base type")
			}

			if t.colWidths[j] < w {
				t.colWidths[j] = w
			}
		}
	}

	return t
}

// Get returns the (i,j)th value.
func (t *Table) Get(i, j int) interface{} {
	return t.body[i][j]
}

// Header returns a copy of the header.
func (t *Table) Header() Header {
	return t.header.Copy()
}

// Import a reader into a table.
func Import(r io.Reader, tableName string, fltFmt FltFmt, fltPrec FltPrecFmt) (*Table, error) {
	// return ImportCSV(*csv.NewReader(r), tableName, fltFmt, fltPrec)
	var (
		t      = New(tableName, fltFmt, fltPrec)
		_, err = t.ReadFrom(r)
	)

	return t, err
}

// ImportCSV imports a csv file into a new table.
func ImportCSV(r csv.Reader, tableName string, fltFmt FltFmt, fltPrec FltPrecFmt) (*Table, error) {
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var (
		t = New(tableName, fltFmt, fltPrec)
		n = len(lines)
	)

	if n == 0 {
		return t, nil
	}

	t.SetHeader(lines[0])
	if n == 1 {
		return t, nil // No body
	}

	for _, line := range lines[1:] {
		r := make(Row, 0, len(line))
		for _, s := range line {
			r = append(r, parse(s))
		}

		t.AppendRow(r)
	}

	return t, nil
}

// minColBaseType returns the smallest base type found in column j.
func (t *Table) minColBaseType(j int) baseType {
	var (
		min = integerType
		bt  baseType
	)

	for _, r := range t.body {
		if bt = baseTypeOf(r[j]); bt < min {
			min = bt
		}
	}

	return min
}

// Read ...
func (t *Table) Read(b []byte) (int, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 0))
	if _, err := t.WriteTo(buf); err != nil {
		return 0, err
	}

	return buf.Read(b)
}

// ReadFrom a reader. The reader should delimit rows with '\n' and
// columns with ','. The first row (the header) will be compared for
// equality against the table's header. For example, given a table
// with a header ["index", "value"], the following buffer contents
// will be read into the table.
// 	buf := bytes.NewBuffer([]byte{})
// 	buf.WriteString("index,value\n0,a\n1,b\n2,c")
//	n, err := table.ReadFrom(buf)
func (t *Table) ReadFrom(r io.Reader) (int64, error) {
	var (
		bts, err = ioutil.ReadAll(r)
		m        = int64(len(bts))
	)

	if err != nil {
		return m, err
	}

	_, err = t.Write(bts)
	return m, err
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

// RemoveRow from a table.
func (t *Table) RemoveRow(i int) Row {
	r := t.body[i]
	if i+1 == t.rows {
		t.body = t.body[:i]
	} else {
		t.body = append(t.body[:i], t.body[i+1:]...)
	}

	t.rows--
	t.Clean()
	return r
}

// Row returns a copy of a row.
func (t *Table) Row(i int) Row {
	return t.body[i].Copy()
}

// Set the (i,j)th cell to a given value.
func (t *Table) Set(v interface{}, i, j int) *Table {
	t.body[i][j] = v
	return t
}

// SetBody ...
func (t *Table) SetBody(b Body) *Table {
	for 0 < t.rows {
		t.RemoveRow(t.rows - 1)
	}

	return t.AppendRows(b...)
}

// SetColHeader to a given value.
func (t *Table) SetColHeader(columnHeader string, i int) *Table {
	t.header[i] = strings.TrimSpace(columnHeader)
	return t
}

// setColSize to a given size n. Empty strings will be appended.
func (t *Table) setColSize(n int) *Table {
	t.columns = n
	for len(t.header) < n {
		t.header = append(t.header, "")
	}

	for n < len(t.header) {
		t.header = append(t.header, "")
	}

	for i := range t.body {
		for len(t.body[i]) < n {
			t.body[i] = append(t.body[i], nil)
		}

		for n < len(t.body[i]) {
			t.body[i] = append(t.body[i], nil)
		}
	}

	for len(t.colBaseTypes) < n {
		t.colBaseTypes = append(t.colBaseTypes, integerType)
	}

	for n < len(t.colBaseTypes) {
		t.colBaseTypes = append(t.colBaseTypes, integerType)
	}

	for len(t.colWidths) < n {
		t.colWidths = append(t.colWidths, 0)
	}

	for n < len(t.colWidths) {
		t.colWidths = append(t.colWidths, 0)
	}

	return t
}

// SetFloatFmt defines how float values are displayed, if any are present.
func (t *Table) SetFloatFmt(f FltFmt) *Table {
	t.floatFmt = f
	return t
}

// SetFloatPrecFmt defines how many digits will be displayed after a decimal
// value, if any are present.
func (t *Table) SetFloatPrecFmt(f FltPrecFmt) *Table {
	t.floatPrecision = f
	return t
}

// SetHeader sets the header field.
func (t *Table) SetHeader(h Header) *Table {
	n := len(h)
	t.setColSize(n)
	for i := 0; i < n; i++ {
		t.header[i] = strings.TrimSpace(h[i])
		if w := len(h[i]); t.colWidths[i] < w { // TODO: should this be t.header[i]?
			t.colWidths[i] = w
		}
	}

	return t
}

// SetMinupdateBaseTypesAndWidths for each table value within the context of its column updateBaseTypesAndWidths.
// That is, this sets the (i,j)th entry to the base type found in each column.
// WARNING: This will wipe out the original data and cannot be undone.
func (t *Table) SetMinupdateBaseTypesAndWidths() *Table {
	for j := 0; j < t.columns; j++ {
		t.colBaseTypes[j] = t.minColBaseType(j)
	}

	for i := 0; i < t.rows; i++ {
		for j := 0; j < t.columns; j++ {
			// Update each (i,j)th value to the jth column base type.
			switch baseTypeOf(t.body[i][j]) {
			case integerType:
				switch t.colBaseTypes[j] {
				case integerType: // Do nothing
				case floatType:
					t.body[i][j] = float64(t.body[i][j].(int)) // Convert to float64
				case stringType:
					t.body[i][j] = strconv.Itoa(t.body[i][j].(int)) // Convert to string
				default:
					panic("unknown base type")
				}
			case floatType:
				switch t.colBaseTypes[j] {
				case integerType: // Do nothing? Data loss if we convert float to int
				case floatType: // Do nothing
				case stringType:
					if x := strconv.FormatFloat(t.body[i][j].(float64), 'f', -1, 64); strings.ContainsRune(x, '.') {
						t.body[i][j] = x
					} else {
						t.body[i][j] = x + ".0"
					}
				default:
					panic("unknown base type")
				}
			case stringType: // Do nothing
			default:
				panic("unknown base type")
			}
		}
	}

	return t
}

// Sort the table's body.
func (t *Table) Sort() *Table {
	sort.Slice(t.body, func(i, j int) bool { return t.body[i].Compare(t.body[j]) < 0 })
	return t
}

// SortOnCol sorts the table body by only comparing the given column index.
func (t *Table) SortOnCol(colIndex int) *Table {
	sort.Slice(t.body, func(i, j int) bool { return t.body[i].CompareAt(t.body[j], colIndex) < 0 })
	return t
}

// String returns a string-representation of a table.
func (t *Table) String() string {
	// Create horizontal line
	var sb strings.Builder
	for _, w := range t.colWidths {
		sb.WriteString("+" + strings.Repeat("-", w))
	}

	sb.WriteString("+\n")
	hLine := sb.String()
	sb.Reset()

	// Write header
	sb.WriteString(hLine)
	for i := 0; i < t.columns; i++ {
		switch t.colBaseTypes[i] {
		case integerType, floatType:
			sb.WriteString("|" + strings.Repeat(" ", t.colWidths[i]-len(t.header[i])) + t.header[i])
		case stringType:
			sb.WriteString("|" + t.header[i] + strings.Repeat(" ", t.colWidths[i]-len(t.header[i])))
		}
	}

	sb.WriteString("|\n" + hLine)

	// Write body
	for i := 0; i < t.rows; i++ {
		for j := 0; j < t.columns; j++ {
			var s string
			switch baseTypeOf(t.body[i][j]) {
			case integerType:
				s = strconv.Itoa(t.body[i][j].(int))
			case floatType:
				s = strconv.FormatFloat(t.body[i][j].(float64), byte(t.floatFmt), int(t.floatPrecision), 64)
			case stringType:
				s = t.body[i][j].(string)
			}

			switch t.colBaseTypes[j] {
			case integerType, floatType:
				sb.WriteString("|" + strings.Repeat(" ", t.colWidths[j]-len(s)) + s)
			case stringType:
				sb.WriteString("|" + s + strings.Repeat(" ", t.colWidths[j]-len(s)))
			}
		}

		sb.WriteString("|\n")
	}

	sb.WriteString(hLine)
	return sb.String()
}

// Strings returns a list of lists-of-strings representing a table.
func (t *Table) Strings() [][]string {
	ss := append(make([][]string, 0, t.rows+1), t.header) // Header is the first row
	for _, r := range t.body {
		ss = append(ss, r.Strings())
	}

	return ss
}

// Swap ...
func (t *Table) Swap(i, j int) {
	t.body.Swap(i, j)
}

// SwapCols ...
func (t *Table) SwapCols(i, j int) {
	t.header.Swap(i, j)
	t.body.SwapCols(i, j)
	t.colBaseTypes[i], t.colBaseTypes[j] = t.colBaseTypes[j], t.colBaseTypes[i]
	t.colWidths[i], t.colWidths[j] = t.colWidths[j], t.colWidths[i]
}

// ExportCSV to a csv writer. Table will be cleaned and set to
// minimum updateBaseTypesAndWidths.
func (t *Table) ExportCSV(w csv.Writer) error {
	return w.WriteAll(t.Clean().SetMinupdateBaseTypesAndWidths().Strings())
}

// Write bytes to a table.
func (t *Table) Write(b []byte) (int, error) {
	m := len(b)
	if m == 0 {
		return m, nil
	}

	lines := toLines(b)
	if t.header.isEmpty() {
		t.header = NewHeader(lines[0]...)
		if len(lines) < 2 {
			return m, nil
		}
		lines = lines[1:]
	} else if t.header.Compare(NewHeader(lines[0]...)) == 0 {
		if len(lines) < 2 {
			return m, nil
		}
		lines = lines[1:]
	}

	for i := range lines {
		r := make(Row, 0, len(lines[i]))
		for j := range lines[i] {
			r = append(r, parse(lines[i][j]))
		}

		t.AppendRow(r)
	}

	return m, nil
}

// WriteTo to a writer.
func (t *Table) WriteTo(w io.Writer) (int64, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := t.ExportCSV(*csv.NewWriter(buf)); err != nil {
		return 0, err
	}

	n, err := w.Write(buf.Bytes())
	return int64(n), err
}
