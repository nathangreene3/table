package main

import "fmt"

func main() {
	t := newTable("Table 1", []alignment{1, 0, 1}, row{"Col A", "   Col B   ", "Col C"}, row{"hello world", "Hi", "Goodbye"})
	t.addRow(row{"5", "6", "7"})
	t.addColumn("Col D", col{"4", "8"}, 1)
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
