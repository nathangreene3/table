package table

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	"github.com/nathangreene3/math"
)

func TestImportExportCSV(t *testing.T) {
	inFile, err := os.Open("test0.csv")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer inFile.Close()

	table, err := ImportFromCSV(csv.NewReader(inFile), "star wars", FltFmtNoExp, 3)
	if err != nil {
		t.Fatalf("%v", err)
	}

	outFile, err := os.OpenFile("test1.csv", os.O_WRONLY, os.ModeAppend)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer outFile.Close()

	if err = table.ExportToCSV(csv.NewWriter(outFile)); err != nil {
		t.Fatalf("%v", err)
	}

	// t.Fatalf("\n%s", table.String())
}

func TestTable(t *testing.T) {
	var (
		left, right    int
		xFacts, yFacts string
		n              = 1 << 10
		tbl            = New("", FltFmtNoExp, 0, n, 6)
	)

	tbl.SetHeader(Header{"x", "y", "(x^2+x)/2", "y^2", "Facts of x", "Facts of y"})
	for x := 0; x < n; x++ {
		left = x * (x + 1) >> 1
		if x == 0 {
			xFacts = ""
		} else {
			xFacts = fmt.Sprint(math.Factor(x))
		}

		for y := 0; y < n; y++ {
			right = y * y
			if y == 0 {
				yFacts = ""
			} else {
				yFacts = fmt.Sprint(math.Factor(y))
			}

			if left == right {
				tbl.AppendRow(Row{x, y, left, right, xFacts, yFacts})
			}
		}
	}

	// t.Fatalf("\n%s\n", tbl.String())
}
