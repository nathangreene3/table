package table2

import (
	"fmt"
	"testing"
)

func TestTable(t *testing.T) {
	{
		tbl := New(NewHeader("Integers", "Floats")).Append(
			NewRow(0, 0.0),
		).Append(
			NewRow(-1, -1.1),
			NewRow(4, -4.4),
			NewRow(-3, 3.3),
		)

		r4 := tbl.Remove(2)
		tbl.Insert(2, NewRow(2, 0.0)).Append(r4).Set(2, 1, 2.2)

		fmt.Println(tbl.String())
		fmt.Println(tbl.Fmts())
		fmt.Println(tbl.InsertCol(0, "Strings", NewCol("zero", "one", "two", "three", "four")).String())
	}

	t.Fatal()
}
