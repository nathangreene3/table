package table

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

// A Table holds rows/columns of data.
type Table struct {
	Name           string
	header         Header
	body           []Row
	rows, columns  int
	colBaseTypes   []baseType
	colWidths      []int
	floatFmt       byte
	floatPrecision int
}

// FloatFormat defines formatting float values.
type FloatFormat byte

const (
	// s="FormatFloat func(f float64, fmt byte, prec int, bitSize int) string
	// FormatFloat converts the floating-point number f to a string, according to the format fmt and precision prec. It rounds the result assuming that the original was obtained from a floating-point value of bitSize bits (32 for float32, 64 for float64).

	// The format fmt is one of 'b' (-ddddp±ddd, a binary exponent), 'e' (-d.dddde±dd, a decimal exponent), 'E' (-d.ddddE±dd, a decimal exponent), 'f' (-ddd.dddd, no exponent), 'g' ('e' for large exponents, 'f' otherwise), or 'G' ('E' for large exponents, 'f' otherwise).

	// The precision prec controls the number of digits (excluding the exponent) printed by the 'e', 'E', 'f', 'g', and 'G' formats. For 'e', 'E', and 'f' it is the number of digits after the decimal point. For 'g' and 'G' it is the maximum number of significant digits (trailing zeros are removed). The special precision -1 uses the smallest number of digits necessary such that ParseFloat will return f exactly."

	// FltFmtBinExp formats floats as a binary exponent value -dddp±ddd.
	FltFmtBinExp FloatFormat = 'b'
	// FltFmtDecExp formats floats as a decimal exponent value (scientific notation) -d.ddde±ddd.
	FltFmtDecExp FloatFormat = 'e'
	// FltFmtNoExp formats floats as a decimal value -ddd.ddd.
	FltFmtNoExp FloatFormat = 'f'
	// FltFmtLrgExp formats floats as a large exponent value (scientific notation) -d.ddde±ddd.
	FltFmtLrgExp FloatFormat = 'g'
)

// New returns an empty table.
func New(name string, floatFmt byte, floatPrec, maxRows, maxColumns int) *Table {
	return &Table{
		Name:           name,
		header:         make(Header, 0, maxColumns),
		body:           make([]Row, 0, maxRows),
		colBaseTypes:   make([]baseType, 0, maxColumns),
		colWidths:      make([]int, 0, maxColumns),
		floatFmt:       floatFmt,
		floatPrecision: floatPrec,
	}
}

// Dimensions returns the number of rows and columns of a table.
func (t *Table) Dimensions() (int, int) {
	return t.rows, t.columns
}

// GetHeader returns a copy of the header.
func (t *Table) GetHeader() Header {
	return t.header.Copy()
}

// SetHeader sets the header field.
func (t *Table) SetHeader(h Header) {
	n := len(h)
	t.setColumns(n)

	var w int
	for i := 0; i < n; i++ {
		t.header[i] = strings.TrimSpace(h[i])
		w = len(h[i])
		if t.colWidths[i] < w {
			t.colWidths[i] = w
		}
	}
}

// GetColumnHeader at a given index.
func (t *Table) GetColumnHeader(i int) string {
	return t.header[i]
}

// SetColumnHeader to a given value.
func (t *Table) SetColumnHeader(columnHeader string, i int) {
	t.header[i] = strings.TrimSpace(columnHeader)
}

// Get returns the (i,j)th value.
func (t *Table) Get(i, j int) interface{} {
	return t.body[i][j]
}

// Set the (i,j)th cell to a given value.
func (t *Table) Set(v interface{}, i, j int) {
	t.body[i][j] = v
}

// GetRow returns a copy of a row.
func (t *Table) GetRow(i int) Row {
	return t.body[i].Copy()
}

// AppendRow to a table.
func (t *Table) AppendRow(r Row) {
	n := len(r)
	t.setColumns(n)
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
				w = len(strconv.FormatFloat(r[i].(float64), t.floatFmt, t.floatPrecision, 64))
			case stringType:
				t.colBaseTypes[i] = stringType
				w = len(r[i].(string))
			default:
				panic("unknown type")
			}
		case floatType:
			switch baseTypeOf(r[i]) {
			case integerType:
				w = len(strconv.FormatFloat(float64(r[i].(int)), t.floatFmt, t.floatPrecision, 64))
			case floatType:
				w = len(strconv.FormatFloat(r[i].(float64), t.floatFmt, t.floatPrecision, 64))
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
				w = len(strconv.FormatFloat(r[i].(float64), t.floatFmt, t.floatPrecision, 64))
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

// GetColumn returns the column header and a copy of the column at a given index.
func (t *Table) GetColumn(i int) (string, Column) {
	c := make(Column, 0, len(t.body))
	for j := range t.body {
		c = append(c, t.body[j][i])
	}

	return t.header[i], c
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

// Format a table.
func (t *Table) Format() {
	for i := 0; i < t.columns; i++ {
		t.colBaseTypes[i] = integerType
		t.colWidths[i] = len(t.header[i])
	}

	var (
		bt baseType
		w  int
	)

	for i := 0; i < t.rows; i++ {
		for j := 0; j < t.columns; j++ {
			if bt = baseTypeOf(t.body[i][j]); bt < t.colBaseTypes[j] {
				t.colBaseTypes[j] = bt
			}

			switch bt {
			case integerType:
				w = len(strconv.Itoa(t.body[i][j].(int)))
			case floatType:
				w = len(strconv.FormatFloat(t.body[i][j].(float64), t.floatFmt, t.floatPrecision, 64))
			case stringType:
				w = len(t.body[i][j].(string))
			default:
				panic("unknown base type")
			}

			if t.colWidths[j] < w {
				t.colWidths[j] = w
			}
		}
	}
}

// SetMinFormat for each table value within the context of its column format.
func (t *Table) SetMinFormat() {
	var (
		bt baseType
		x  string
	)
	for i := 0; i < t.rows; i++ {
		for j := 0; j < t.columns; j++ {
			if bt = baseTypeOf(t.body[i][j]); bt < t.colBaseTypes[j] {
				t.colBaseTypes[j] = bt
			}

			switch bt {
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
					x = strconv.FormatFloat(t.body[i][j].(float64), 'f', -1, 64)
					if strings.ContainsRune(x, '.') {
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

// setColumns to a given size n.
func (t *Table) setColumns(n int) {
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

// ImportFromCSV imports a csv file into a table and returns it.
func ImportFromCSV(reader *csv.Reader, tableName string, fltFmt byte, fltPrec int) (*Table, error) {
	t := New(tableName, fltFmt, fltPrec, 0, 0)

	// Header
	line, err := reader.Read()
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
		return t, nil
	}

	t.SetHeader(line)

	// Body
	for {
		line, err = reader.Read()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return t, nil
		}

		r := make(Row, 0, len(line))
		for _, s := range line {
			r = append(r, toBaseType(s))
		}

		t.AppendRow(r)
	}
}

// ExportToCSV to a given path. Table will be cleaned and set to minimum format.
func (t *Table) ExportToCSV(writer *csv.Writer) error {
	t.Clean()
	t.SetMinFormat()

	writer.Write([]string(t.header))
	for _, r := range t.body {
		writer.Write(r.Strings())
	}

	writer.Flush()
	return writer.Error()
}

// Copy a table.
func (t *Table) Copy() *Table {
	cpy := New(t.Name, t.floatFmt, t.floatPrecision, t.rows, t.columns)
	cpy.SetHeader(t.header)
	for i := 0; i < t.rows; i++ {
		cpy.AppendRow(t.body[i].Copy())
	}

	return cpy
}

// String returns a string-representation of a table.
func (t *Table) String() string {
	// Create horizontal line
	sb := strings.Builder{}
	for i := range t.colWidths {
		sb.WriteString("+" + strings.Repeat("-", t.colWidths[i]))
	}

	sb.WriteString("+\n")
	hLine := sb.String()
	sb.Reset()

	// Write header
	sb.WriteString(hLine)
	for i := 0; i < t.columns; i++ {
		switch t.colBaseTypes[i] {
		case integerType:
			sb.WriteString("|" + t.header[i] + strings.Repeat(" ", t.colWidths[i]-len(t.header[i])))
		case floatType:
			sb.WriteString("|" + t.header[i] + strings.Repeat(" ", t.colWidths[i]-len(t.header[i])))
		case stringType:
			sb.WriteString("|" + strings.Repeat(" ", t.colWidths[i]-len(t.header[i])) + t.header[i])
		}
	}

	sb.WriteString("|\n")
	sb.WriteString(hLine)

	// Write body
	var s string
	for i := 0; i < t.rows; i++ {
		for j := 0; j < t.columns; j++ {
			switch baseTypeOf(t.body[i][j]) {
			case integerType:
				s = strconv.Itoa(t.body[i][j].(int))
			case floatType:
				s = strconv.FormatFloat(t.body[i][j].(float64), t.floatFmt, t.floatPrecision, 64)
			case stringType:
				s = t.body[i][j].(string)
			}

			switch t.colBaseTypes[j] {
			case integerType:
				sb.WriteString("|" + strings.Repeat(" ", t.colWidths[j]-len(s)) + s)
			case floatType:
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
