package main

import "fmt"

func main() {
	t := newTable(row{"Col A", "   Col B   ", "Col C"})
	t.addRow(row{"1", "2", "3"})
	t.addColumn("Col D", col{"4", "8"})
	fmt.Println(t.String())
}

func maxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}
