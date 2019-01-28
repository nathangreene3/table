package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

func main() {
	t := newTable("Random Numbers")
	t.addColumn("Index", col{}, alignLeft)
	t.addColumn("Value", col{}, alignRight)
	for i := 0; i < 5; i++ {
		t.addRow(row{strconv.Itoa(i + 1), strconv.Itoa(rand.Intn(10))})
	}
	fmt.Println(t.String())
}
