package main

import "fmt"

func main() {
	t := newTable("Table 1", []alignment{1, 1, 1}, row{"Col A", "   Col B   ", "Col C"}, row{"hello world", "Hi", "Goodbye"})
	t.addRow(row{"5", "6", "7"})
	t.addColumn("Col D", col{"4", "8"}, 1)
	// t.removeRow(0)
	// t.removeColumn(1)
	t.setCell(4, 4, "  new cell  ")
	t.setHeader(append(t.header[:t.width-1], "Col E"))
	t.setAlignment(append(t.align[:t.width-1], alignRight))
	fmt.Println(t.String())
}

func maxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}
