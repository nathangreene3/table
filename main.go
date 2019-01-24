package main

import "fmt"

func main() {
	t := newTable("Table A", nil, nil, nil)
	t.addColumn("Index", col{}, alignRight)
	t.addColumn("First Name", col{}, alignLeft)
	t.addColumn("Last Name", col{}, alignLeft)
	t.addRow(row{"1", "Nathan", "Greene"})
	t.addRow(row{"2", "Sarah", "Cronk"})
	t.addRow(row{"3", "Grace", "Greene"})
	fmt.Println(t.String())
	fmt.Println(t.info(), t.maxColWidth)
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
