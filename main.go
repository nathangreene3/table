package main

import "fmt"

func main() {
	t := newTable("Table A")
	t.addColumn("Index", col{}, alignRight)
	t.addColumn("First Name", col{}, alignCenter)
	t.addColumn("Last Name", col{}, alignLeft)
	t.addRow(row{"1", "Nathan", "Greene"})
	t.addRow(row{"2", "Sarah", "Cronk"})
	t.addRow(row{"3", "Grace", "Greene"})
	t.addRow(row{"4", "Jim Jimmy Jim Jim Jim Jim Jim", "", "0"})
	t.setColHeader(t.width-1, "Age", alignRight)
	t.setCell(0, 3, "28")
	t.setCell(1, 3, "25")
	t.setCell(2, 3, "20")
	// t.removeColumn(1)
	fmt.Println(t.String())
	// fmt.Println(t.info(), t.maxColWidth)
}

func maxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}
