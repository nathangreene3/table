package main

import (
	"strings"
)

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
)

// String returns a formated representation of a table.
func (t *table) String() string {
	var pad string

	s := make([]string, 0, 2) // Add header and underlines to height

	// Add overline
	s = append(s, strings.Repeat("-", maxInt(t.maxColWidth[0], 1)))
	for i := 1; i < t.width; i++ {
		s = append(s, strings.Repeat("-", maxInt(t.maxColWidth[i], 1)+1))
	}

	// Add name
	s = append(s, "\n"+t.name+"\n")

	// Add header
	pad = strings.Repeat(" ", t.maxColWidth[0]-len(t.header[0]))
	switch t.align[0] {
	case alignRight:
		s = append(s, pad+t.header[0])
	default:
		s = append(s, t.header[0]+pad)
	}
	for i := 1; i < t.width; i++ {
		pad = strings.Repeat(" ", t.maxColWidth[i]-len(t.header[i]))
		switch t.align[i] {
		case alignRight:
			s = append(s, " "+pad+t.header[i])
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
		pad = strings.Repeat(" ", t.maxColWidth[0]-len(t.body[i][0]))
		switch t.align[0] {
		case alignRight:
			s = append(s, "\n"+pad+t.body[i][0])
		default:
			s = append(s, "\n"+t.body[i][0]+pad)
		}
		for j := 1; j < t.width; j++ {
			pad = strings.Repeat(" ", t.maxColWidth[j]-len(t.body[i][j]))
			switch t.align[j] {
			case alignRight:
				s = append(s, " "+pad+t.body[i][j])
			default:
				s = append(s, " "+t.body[i][j]+pad)
			}
		}
	}

	// Add underline
	s = append(s, "\n"+strings.Repeat("-", maxInt(t.maxColWidth[0], 1)))
	for i := 1; i < t.width; i++ {
		s = append(s, strings.Repeat("-", maxInt(t.maxColWidth[i], 1)+1))
	}

	return strings.Join(s, "")
}

// newTable returns a table with the header set and an empty body. The height is set to zero and the width is set to the length of header.
func newTable(name string, align []alignment, header row, body ...row) *table {
	t := &table{
		name:        strings.TrimSpace(name),
		header:      trimRow(header),
		body:        make([]row, 0, len(body)),
		height:      0,
		width:       len(header),
		maxColWidth: make([]int, len(header)),
		align:       make([]alignment, len(header)),
	}

	for i := range align {
		t.align[i] = align[i]
	}

	t.updateMaxColWidths()

	for i := range body {
		t.addRow(body[i])
	}

	return t
}

func (t *table) info() map[string]interface{} {
	m := make(map[string]interface{})
	m["name"] = t.name
	m["header"] = t.header
	m["rows"] = t.height
	m["columns"] = t.width
	m["alignment"] = t.align
	return m
}

func (t *table) setName(name string) {
	t.name = strings.TrimSpace(name)
}

// setHeight expands or contracts the table body. If contracting, rows may be removed.
func (t *table) setHeight(height int) {
	if height < t.height {
		t.body = t.body[:height]
		t.height = height
	}

	for t.height < height {
		t.body = append(t.body, make(row, t.width))
		t.height++
	}
}

// setWidth expands or contracts the table. If contracting, columns may be removed.
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

	// Shrink header, body, maximum column width, and align
	if width < len(t.header) {
		t.header = t.header[:width]
	}
	for i := range t.body {
		if width < len(t.body) {
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

// setHeader sets the header. If the new header is smaller than the table width, it is expanded. If it is larger, the table is expanded.
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

// addRow appends a row. If the new row is smaller than the table width, it is expanded. If it is larger, the table is expanded.
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

// addColumn appends a column with its header. If the new column is smaller than the table height, it is expanded. If it is larger, the table is expanded.
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
		t.body = append(t.body[:i], t.body[i+1:]...)
	} else {
		t.body = t.body[:i]
	}
	t.height = len(t.body)

	t.updateMaxColWidths()

	return r
}

// removeColumn removes and returns a column with its header.
func (t *table) removeColumn(i int) (string, col) {
	header := t.header[i]
	c := make(col, 0, t.height)
	if i+1 < t.width {
		for j := range t.body {
			t.header = append(t.header[:i], t.header[i+1:]...)
			c = append(c, t.body[j][i])
			t.body[j] = append(t.body[j][:i], t.body[j][i+1:]...)
		}
	} else {
		for j := range t.body {
			t.header = t.header[:i]
			c = append(c, t.body[j][i])
			t.body[j] = t.body[j][:i]
		}
	}
	t.width = len(t.header)
	return header, c
}

// updateMaxColWidths sets each entry in maxColWidths to the length of the largest entry in each column.
func (t *table) updateMaxColWidths() {
	for i := range t.header {
		t.maxColWidth[i] = len(t.header[i])
	}

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
	for k := range t.body {
		t.body[k][i], t.body[k][j] = t.body[k][j], t.body[k][i]
	}
}

// trimRow returns a deep copy of a row with the leading and trailing whitespace removed from each entry.
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
