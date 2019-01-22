package main

import (
	"strings"
)

type table struct {
	header      row   // Ordered set of descriptors for the corresponding body columns
	body        []row // Ordered set of rows
	maxColWidth []int // Maximum column width; each column is as wide as its largest header or body column member
	height      int   // Number of rows in body
	width       int   // Number of columns in header, body, and maxColWidth
}

type row []string
type col []string

// String returns a formated representation of a table.
func (t *table) String() string {
	s := make([]string, 0, (t.height+2)*t.width) // Add header and underlines to height

	// Add header
	s = append(s, t.header[0])
	for i := 1; i < t.width; i++ {
		s = append(s, strings.Repeat(" ", t.maxColWidth[i-1]-len(t.header[i-1])+1)+t.header[i])
	}
	s = append(s, "\n")

	// Add underlines
	s = append(s, strings.Repeat("-", maxInt(t.maxColWidth[0], 1)))
	for i := 1; i < t.width; i++ {
		s = append(s, " "+strings.Repeat("-", maxInt(t.maxColWidth[i], 1)))
	}
	s = append(s, "\n")

	// Add body
	for i := range t.body {
		s = append(s, t.body[i][0])
		for j := 1; j < t.width; j++ {
			s = append(s, strings.Repeat(" ", t.maxColWidth[j-1]-len(t.body[i][j-1])+1)+t.body[i][j])
		}
		s = append(s, "\n")
	}

	return strings.Join(s, "")
}

// newTable returns a table with the header set and an empty body. The height is set to zero and the width is set to the length of header.
func newTable(header row, body ...row) *table {
	n := len(body)
	t := &table{
		header:      trimRow(header),
		body:        make([]row, 0, n),
		height:      0,
		width:       len(header),
		maxColWidth: make([]int, len(header)),
	}
	t.updateMaxColWidths()

	for i := range body {
		t.addRow(body[i])
	}

	return t
}

func (t *table) setWidth(width int) {
	// Extend header, body, and maximum column width
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

	// Shrink header, body, and maximum column width
	if width < len(t.header) {
		t.header = t.header[:width]
	}
	for i := range t.body {
		if width < len(t.body) {
			t.body[i] = t.body[i][:width]
		}
	}
	for width < len(t.maxColWidth) {
		t.maxColWidth = t.maxColWidth[:width]
	}

	t.width = width
}

func (t *table) setHeader(header row) {
	width := len(header)

	// Extend header, body, and maximum column width
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

func (t *table) addRow(r row) {
	width := len(r)

	// Extend header, body, and maximum column width
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
		if t.maxColWidth[t.width-1] < n {
			t.maxColWidth[t.width-1] = n
		}
	}
}

func (t *table) addColumn(header string, c col) {
	// Extend header
	t.setHeader(append(t.header, header))

	// Extend body
	height := len(c)
	for t.height < height {
		t.addRow(make(row, t.width))
	}

	// Extend new column
	for height < t.height {
		c = append(c, "")
		height++
	}

	// Insert into each row the new column entry
	var n int
	for i := range t.body {
		t.body[i][t.width-1] = c[i]
		n = len(c[i])
		if t.maxColWidth[t.width-1] < n {
			t.maxColWidth[t.width-1] = n
		}
	}
}

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

func (t *table) swapRows(i, j int) {
	t.body[i], t.body[j] = t.body[j], t.body[i]
}

func (t *table) swapCols(i, j int) {
	for k := range t.body {
		t.body[k][i], t.body[k][j] = t.body[k][j], t.body[k][i]
	}
}

func trimRow(r row) row {
	s := make(row, 0, len(r))
	for i := range r {
		s = append(s, strings.TrimSpace(r[i]))
	}
	return s
}

func copyRow(r row) row {
	s := make(row, len(r))
	copy(s, r)
	return s
}
