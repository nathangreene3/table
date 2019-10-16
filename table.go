package table

import (
	"encoding/csv"
	"strconv"
	"strings"
)

// A Table holds rows/columns of data.
type Table struct {
	Name           string
	header         Header
	body           []Row
	colBaseTypes   []baseType
	colWidths      []int
	rows, columns  int
	floatPrecision FltPrecFmt
	floatFmt       FltFmt
}

// FltFmt defines formatting float values.
type FltFmt byte

// FltPrecFmt defines the number of decimal positions displayed for float values.
type FltPrecFmt int

const (
	// FltFmtBinExp formats floats as a binary exponent value -dddp±ddd.
	FltFmtBinExp FltFmt = 'b'
	// FltFmtDecExp formats floats as a decimal exponent value (scientific notation) -d.ddde±ddd.
	FltFmtDecExp FltFmt = 'e'
	// FltFmtNoExp formats floats as a decimal value -ddd.ddd.
	FltFmtNoExp FltFmt = 'f'
	// FltFmtLrgExp formats floats as a large exponent value (scientific notation) -d.ddde±ddd.
	FltFmtLrgExp FltFmt = 'g'
)

// New returns an empty table.
func New(name string, floatFmt FltFmt, floatPrec FltPrecFmt) Table {
	if floatFmt == 0 {
		floatFmt = FltFmtNoExp
	}

	var (
		maxRows = 256
		maxCols = 256
	)

	return Table{
		Name:           name,
		header:         make(Header, 0, maxCols),
		body:           make([]Row, 0, maxRows),
		colBaseTypes:   make([]baseType, 0, maxCols),
		colWidths:      make([]int, 0, maxCols),
		floatPrecision: floatPrec,
		floatFmt:       floatFmt,
	}
}

// AppendColumn to a table.
func (t *Table) AppendColumn(columnHeader string, c Column) {
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
}

// AppendRow to a table.
func (t *Table) AppendRow(r Row) {
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
}

// Clean removes empty rows and columns.
func (t *Table) Clean() {
	// Remove empty rows
	for i := 0; i < t.rows; i++ {
		if t.body[i].isEmpty() {
			t.RemoveRow(i)
		}
	}

	// Remove empty columns
	var isEmpty bool
	for j := 0; j < t.columns; j++ {
		if isEmpty = t.header[j] == ""; !isEmpty {
			continue
		}

		for i := 0; i < t.rows && isEmpty; i++ {
			isEmpty = t.body[i][j] == nil
		}

		if isEmpty {
			t.RemoveColumn(j)
		}
	}

	t.Format()
}

// Column returns the column header and a copy of the column at a given index.
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

// Copy a table.
func (t *Table) Copy() Table {
	cpy := New(t.Name, t.floatFmt, t.floatPrecision)
	cpy.SetHeader(t.header)
	for i := 0; i < t.rows; i++ {
		cpy.AppendRow(t.body[i].Copy())
	}

	return cpy
}

// Dimensions returns the number of rows and columns of a table.
func (t *Table) Dimensions() (int, int) {
	return t.rows, t.columns
}

// Export to a csv writer. Table will be cleaned and set to minimum format.
func (t *Table) Export(writer csv.Writer) error {
	t.Clean()
	t.SetMinFormat()

	writer.Write([]string(t.header))
	for _, r := range t.body {
		writer.Write(r.Strings())
	}

	writer.Flush()
	return writer.Error()
}

// Format a table. This updates each column base type to its weakest base type and updates each column width to the largest each needs to be when formated as a string.
func (t *Table) Format() {
	// Reset each column format as an integer and to be as wide as the column
	// header. Each column base type is used to determine alignment and doesn't
	// affect the formatting of each (i,j)th value.
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
}

// Get returns the (i,j)th value.
func (t *Table) Get(i, j int) interface{} {
	return t.body[i][j]
}

// Header returns a copy of the header.
func (t *Table) Header() Header {
	return t.header.Copy()
}

// Import imports a csv file into a table and returns it.
func Import(reader csv.Reader, tableName string, fltFmt FltFmt, fltPrec FltPrecFmt) (Table, error) {
	var (
		t          = New(tableName, fltFmt, fltPrec)
		lines, err = reader.ReadAll()
		n          = len(lines)
	)

	if err != nil || n == 0 {
		return t, err // No header, no body
	}

	t.SetHeader(lines[0])
	if n == 1 {
		return t, nil // No body
	}

	for _, line := range lines[1:] {
		r := make(Row, 0, len(line))
		for _, s := range line {
			r = append(r, toBaseType(s))
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
	return r
}

// Row returns a copy of a row.
func (t *Table) Row(i int) Row {
	return t.body[i].Copy()
}

// Set the (i,j)th cell to a given value.
func (t *Table) Set(v interface{}, i, j int) {
	t.body[i][j] = v
}

// SetColHeader to a given value.
func (t *Table) SetColHeader(columnHeader string, i int) {
	t.header[i] = strings.TrimSpace(columnHeader)
}

// setColSize to a given size n. Empty strings will be appended.
func (t *Table) setColSize(n int) {
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
}

// SetFloatFmt defines how float values are displayed, if any are present.
func (t *Table) SetFloatFmt(f FltFmt) {
	t.floatFmt = f
}

// SetFloatPrecFmt defines how many digits will be displayed after a decimal
// value, if any are present.
func (t *Table) SetFloatPrecFmt(f FltPrecFmt) {
	t.floatPrecision = f
}

// SetHeader sets the header field.
func (t *Table) SetHeader(h Header) {
	n := len(h)
	t.setColSize(n)

	var w int
	for i := 0; i < n; i++ {
		t.header[i] = strings.TrimSpace(h[i])
		w = len(h[i])
		if t.colWidths[i] < w {
			t.colWidths[i] = w
		}
	}
}

// SetMinFormat for each table value within the context of its column format.
// That is, this sets the (i,j)th entry to the base type found in each column.
// WARNING: This will wipe out the original data and cannot be undone.
func (t *Table) SetMinFormat() {
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
	var s string
	for i := 0; i < t.rows; i++ {
		for j := 0; j < t.columns; j++ {
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
