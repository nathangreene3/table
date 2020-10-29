package table2

import (
	"fmt"
	"testing"
)

func TestTable(t *testing.T) {
	{
		var (
			h   = Header{"Integers", "Floats", "Strings"}
			tbl = New(h).Append(Row{0, 0.0, "zero"})
		)

		tbl.Append(
			Row{1, 1.1, "one"},
			Row{2, 2.2, "two"},
			Row{3, 3.3, "three"},
		)

		fmt.Println(tbl.String())
		fmt.Println(tbl.Fmt())
	}

	t.Fatal()
}
