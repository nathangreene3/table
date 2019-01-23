# GoTable

## Usage

```Go
package main

import (
    "fmt",
    "github.com/nathangreene3/GoTable",
)

func main() {
	t := newTable("Table 1", []alignment{1, 1, 1}, row{"Col A", "   Col B   ", "Col C"}, row{"hello world", "Hi", "Goodbye"})
	t.addRow(row{"5", "6", "7"})
	t.addColumn("Col D", col{"4", "8"}, alignRight)
	t.removeRow(0)
	t.removeColumn(1)
	t.setCell(4, 4, "  new cell  ")
	t.setHeader(append(t.header[:t.width-1], "Col E"))
	t.setAlignment(append(t.align[:t.width-1], alignRight))
	t.setCell(5, 5, "123")
	t.setColHeader(5, "Col E", alignRight)
	fmt.Println(t.String())
}
```