package table2

import (
	"fmt"
	"testing"
)

func TestTable(t *testing.T) {
	{
		tbl := New(NewHeader("Integers", "Floats", "Strings")).Append(
			NewRow(0, 0.0, "zero"),
		).Append(
			NewRow(1, 1.1, "one"),
			NewRow(4, 4.4, "four"),
			NewRow(3, 3.3, "three"),
		)

		r4 := tbl.Remove(2)
		tbl.Insert(2, NewRow(2, 0.0, "two")).Append(r4).Set(2, 1, 2.2)

		fmt.Println(tbl.String())
		fmt.Println(tbl.Fmts())
	}

	t.Fatal()
}
