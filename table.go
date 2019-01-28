package main

import (
	"strings"
)

// A table is a set of cells containing string entries.
type table struct {
	name        string      // Table name
	header      row         // Ordered set of descriptors for the corresponding body columns
	body        []row       // Ordered set of rows
	maxColWidth []int       // Maximum column width; each column is as wide as its largest header or body column member
	height      int         // Number of rows in body
	width       int         // Number of columns in header, body, and maxColWidth
	align       []alignment // TODO: add align left and right

}

type row []string
type col []string
type alignment int

const (
	alignLeft alignment = iota
	alignRight
	alignCenter
)

// String returns a formated representation of a table.
func (t *table) String() string {
	// Empty table body condition
	var line string
	if t.width == 0 {
		line = strings.Repeat("-", maxInt(len(t.name), 1))
		return line + "\n" + t.name + "\n" + line
	}

	s := make([]string, 0, t.length()) // String to build and return

	// Create over/underline
	var n int // Temporary storage
	n = t.width - 1
	for i := range t.maxColWidth {
		n += maxInt(t.maxColWidth[i], 1)
	}
	line = strings.Repeat("-", maxInt(n, len(t.name)))

	// Add overline and name
	s = append(s, line+"\n"+t.name+"\n")

	// Add header
	var pad string
	if t.maxColWidth[0] == 0 {
		pad = " "
	} else {
		pad = strings.Repeat(" ", t.maxColWidth[0]-len(t.header[0]))
	}

	switch t.align[0] {
	case alignRight:
		s = append(s, pad+t.header[0])
	case alignCenter:
		n = len(pad) / 2
		s = append(s, pad[:n]+t.header[0]+pad[n:])
	case alignLeft:
		fallthrough
	default:
		s = append(s, t.header[0]+pad)
	}

	for i := 1; i < t.width; i++ {
		if t.maxColWidth[i] == 0 {
			pad = " "
		} else {
			pad = strings.Repeat(" ", t.maxColWidth[i]-len(t.header[i]))
		}

		switch t.align[i] {
		case alignRight:
			s = append(s, " "+pad+t.header[i])
		case alignCenter:
			n = len(pad) / 2
			s = append(s, " "+pad[:n]+t.header[i]+pad[n:])
		case alignLeft:
			fallthrough
		default:
			s = append(s, " "+t.header[i]+pad)
		}
	}

	// Add underlines
	s = append(s, "\n"+strings.Repeat("-", maxInt(t.maxColWidth[0], 1)))
	for i := 1; i < t.width; i++ {
		s = append(s, " "+strings.Repeat("-", maxInt(t.maxColWidth[i], 1)))
	}

	// Add body
	for i := range t.body {
		if t.maxColWidth[0] == 0 {
			pad = " "
		} else {
			pad = strings.Repeat(" ", t.maxColWidth[0]-len(t.body[i][0]))
		}

		switch t.align[0] {
		case alignRight:
			s = append(s, "\n"+pad+t.body[i][0])
		case alignCenter:
			n = len(pad) / 2
			s = append(s, "\n"+pad[:n]+t.body[i][0]+pad[:n])
		case alignLeft:
			fallthrough
		default:
			s = append(s, "\n"+t.body[i][0]+pad)
		}

		for j := 1; j < t.width; j++ {
			if t.maxColWidth[j] == 0 {
				pad = " "
			} else {
				pad = strings.Repeat(" ", t.maxColWidth[j]-len(t.body[i][j]))
			}

			switch t.align[j] {
			case alignRight:
				s = append(s, " "+pad+t.body[i][j])
			case alignCenter:
				n = len(pad) / 2
				s = append(s, " "+pad[:n]+t.body[i][j]+pad[n:])
			case alignLeft:
				fallthrough
			default:
				s = append(s, " "+t.body[i][j]+pad)
			}
		}
	}

	// Add underline
	s = append(s, "\n"+line)

	return strings.Join(s, "")
}

// newTable returns a named, empty table.
func newTable(name string) *table {
	return &table{
		name:        strings.TrimSpace(name),
		header:      make(row, 0, 8),
		body:        make([]row, 0, 256),
		height:      0,
		width:       0,
		maxColWidth: make([]int, 0, 8),
		align:       make([]alignment, 0, 8),
	}
}

// info returns a map of the table fields.
func (t *table) info() map[string]interface{} {
	m := make(map[string]interface{})
	m["name"] = t.name
	m["header"] = t.header
	m["rows"] = t.height
	m["columns"] = t.width
	m["alignment"] = t.align
	return m
}

// setCell sets a cell to a given value. If the indices are greater than
// the width and height, then the table is extended with empty rows and
// columns left aligned by default.
func (t *table) setCell(i, j int, s string) {
	if t.height <= i {
		t.setHeight(i + 1)
	}
	if t.width <= j {
		t.setWidth(j + 1)
	}
	t.body[i][j] = strings.TrimSpace(s)
	t.updateMaxColWidths()
}

// setName sets the name of the table.
func (t *table) setName(name string) {
	t.name = strings.TrimSpace(name)
}

// setHeight expands or contracts the table. If expanding, empty rows may
// be added. If contracting, rows may be removed.
func (t *table) setHeight(height int) {
	// Contract body
	if height < t.height {
		t.body = t.body[:height]
		t.height = height
	}

	// Extend body
	for t.height < height {
		t.addRow(make(row, t.width))
	}
}

// setWidth expands or contracts the table. If expanding, empty columns
// may be added. If contracting, columns may be removed.
func (t *table) setWidth(width int) {
	// Extend header, body, maximum column width, and align
	for len(t.header) < width {
		t.header = append(t.header, "")
	}
	for i := range t.body {
		for len(t.body[i]) < width {
			t.body[i] = append(t.body[i], "")
		}
	}
	for len(t.maxColWidth) < width {
		t.maxColWidth = append(t.maxColWidth, 0)
	}
	for len(t.align) < width {
		t.align = append(t.align, alignLeft)
	}

	// Contract header, body, maximum column width, and align
	if width < len(t.header) {
		t.header = t.header[:width]
	}
	for i := range t.body {
		if width < len(t.body[i]) {
			t.body[i] = t.body[i][:width]
		}
	}
	if width < len(t.maxColWidth) {
		t.maxColWidth = t.maxColWidth[:width]
	}
	if width < len(t.align) {
		t.align = t.align[:width]
	}

	t.width = width
}

// setAlignment sets the table alignment for each column. If contracting,
// columns may be removed.
func (t *table) setAlignment(align []alignment) {
	width := len(align)

	// Extend header, body, maximum column width, and align
	if t.width < width {
		t.setWidth(width)
	}

	// Extend new align
	for width < t.width {
		align = append(align, alignLeft)
		width++
	}

	copy(t.align, align)
}

// setHeader sets the header. If the new header is smaller than the table
// width, it is expanded. If it is larger, the table is expanded.
func (t *table) setHeader(header row) {
	width := len(header)

	// Extend header, body, maximum column width, and align
	if t.width < width {
		t.setWidth(width)
	}

	// Extend new header
	for width < t.width {
		header = append(header, "")
		width++
	}

	copy(t.header, trimRow(header))

	// Update maximum column widths
	var n int
	for i := range t.header {
		n = len(t.header[i])
		if t.maxColWidth[i] < n {
			t.maxColWidth[i] = n
		}
	}
}

// setColHeader sets a column header name and aligns it. If the column
// index is wider than the table width, then the table is extended.
func (t *table) setColHeader(i int, header string, align alignment) {
	if t.width <= i {
		t.setWidth(i + 1)
	}

	t.header[i] = strings.TrimSpace(header)
	t.updateMaxColWidths()
	t.align[i] = align
}

// addRow appends a row to the table. If the new row is smaller than the
// table width, it is expanded. If it is larger, the table is expanded.
func (t *table) addRow(r row) {
	width := len(r)

	// Extend header, body, maximum column width, and align
	if t.width < width {
		t.setWidth(width)
	}

	// Extend new row
	for width < t.width {
		r = append(r, "")
		width++
	}

	t.body = append(t.body, trimRow(r))
	t.height = len(t.body)

	// Update maximum column widths
	var n int
	for i := range r {
		n = len(r[i])
		if t.maxColWidth[i] < n {
			t.maxColWidth[i] = n
		}
	}
}

// addColumn appends a column with its header to the table. If the new
// column is smaller than the table height, it is expanded. If it is
// larger, the table is expanded.
func (t *table) addColumn(header string, c col, align alignment) {
	height := len(c)

	// Extend table
	t.setWidth(t.width + 1)

	// Set header and alignment
	t.header[t.width-1] = strings.TrimSpace(header)
	t.align[t.width-1] = align

	// Update maximum column width
	n := len(t.header[t.width-1])
	if t.maxColWidth[t.width-1] < n {
		t.maxColWidth[t.width-1] = n
	}

	// Extend body
	if t.height < height {
		t.setHeight(height)
	}

	// Extend new column
	for height < t.height {
		c = append(c, "")
		height++
	}

	// Insert into each row the new column entry
	for i := range t.body {
		t.body[i][t.width-1] = c[i]
		n = len(c[i])
		if t.maxColWidth[t.width-1] < n {
			t.maxColWidth[t.width-1] = n
		}
	}
}

// removeRow removes and returns a row.
func (t *table) removeRow(i int) row {
	// Extract row to return
	r := t.body[i]
	if i+1 < t.height {
		// Remove row from somewhere before the last row
		t.body = append(t.body[:i], t.body[i+1:]...)
	} else {
		// Remove last row
		t.body = t.body[:i]
	}
	t.height--
	t.updateMaxColWidths()
	return r
}

// removeColumn removes and returns a column with its header.
func (t *table) removeColumn(i int) (string, col) {
	// Copy column info
	header := t.header[i]       // Column header to return
	c := make(col, 0, t.height) // Column to return
	for j := range t.body {
		c = append(c, t.body[j][i])
	}

	// Remove column
	t.header = append(t.header[:i], t.header[i+1:]...)
	t.maxColWidth = append(t.maxColWidth[:i], t.maxColWidth[i+1:]...)
	t.align = append(t.align[:i], t.align[i+1:]...)
	for j := range t.body {
		t.body[j] = append(t.body[j][:i], t.body[j][i+1:]...)
	}
	t.width--

	return header, c
}

// updateMaxColWidths sets each entry in maxColWidths to the length of
// the largest entry in each column.
func (t *table) updateMaxColWidths() {
	// Set maximum column width to length of each header
	for i := range t.header {
		t.maxColWidth[i] = len(t.header[i])
	}

	// Update maximum column width if each entry is larger
	var n int
	for i := range t.body {
		for j := range t.body[i] {
			n = len(t.body[i][j])
			if t.maxColWidth[j] < n {
				t.maxColWidth[j] = n
			}
		}
	}
}

// swapRows swaps two rows.
func (t *table) swapRows(i, j int) {
	t.body[i], t.body[j] = t.body[j], t.body[i]
}

// swapCols swaps two columns.
func (t *table) swapCols(i, j int) {
	t.header[i], t.header[j] = t.header[j], t.header[i]
	t.maxColWidth[i], t.maxColWidth[j] = t.maxColWidth[j], t.maxColWidth[i]
	t.align[i], t.align[j] = t.align[j], t.align[i]
	for k := range t.body {
		t.body[k][i], t.body[k][j] = t.body[k][j], t.body[k][i]
	}
}

// trimRow returns a deep copy of a row with the leading and trailing
// whitespace removed from each entry.
func trimRow(r row) row {
	s := make(row, 0, len(r))
	for i := range r {
		s = append(s, strings.TrimSpace(r[i]))
	}
	return s
}

// copyRow returns a deep copy of a row.
func copyRow(r row) row {
	s := make(row, len(r))
	copy(s, r)
	return s
}

func (t *table) length() int {
	// Empty table
	if t.width == 0 {
		return 3*len(t.name) + 2
	}

	// Non-empty table
	var s int
	for i := range t.maxColWidth {
		s += t.maxColWidth[i] + 1 // +1 refers to spaces and \n
	}
	return s * (t.height + 5)
}

// maxInt returns the maximum of two integers.
func maxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// minInt returns the minimum of two integers.
func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}
