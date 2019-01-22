package main

import "fmt"

func main() {
	t := newTable(row{"Col A", "   Col B   ", "Col C"}, row{"1", "2", "3"})
	t.addRow(row{"5", "6", "7"})
	t.addColumn("Col D", col{"4", "8"})
	// t.removeRow(0)
	// t.removeColumn(1)
	fmt.Println(t.String())
}

func maxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}
