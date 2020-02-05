package table

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/nathangreene3/math"
)

func TestImportExport(t *testing.T) {
	inFile, err := os.Open("test0.csv")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer inFile.Close()

	table, err := Import(inFile, "star wars", FltFmtNoExp, 3)
	if err != nil {
		t.Fatalf("%v", err)
	}

	outFile, err := os.OpenFile("test1.csv", os.O_WRONLY, os.ModeAppend)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer outFile.Close()

	if err = table.Export(outFile); err != nil {
		t.Fatalf("%v", err)
	}

	// t.Fatalf("\n%s", table.String())
}

func TestTable(t *testing.T) {
	type factPow struct {
		factor, power int
	}

	factorsStr := func(n int) string {
		if n < 1 {
			return ""
		}

		var (
			factors  = math.Factor(n)
			factPows = make([]factPow, 0, len(factors))
		)

		for factor, power := range factors {
			factPows = append(factPows, factPow{factor: factor, power: power})
		}

		sort.Slice(
			factPows,
			func(i, j int) bool {
				switch {
				case factPows[i].factor < factPows[j].factor:
					return true
				case factPows[i].factor == factPows[j].factor:
					return factPows[i].power < factPows[j].power
				default:
					return false
				}
			},
		)

		s := make([]string, 0, len(factors))
		for _, fs := range factPows {
			s = append(s, fmt.Sprintf("%d^%d", fs.factor, fs.power))
		}

		return strings.Join(s, " * ")
	}

	var (
		index int
		n     = 1 << 16
		tbl   = New("Squared-Triangle Numbers", FltFmtNoExp, 0).SetHeader(NewHeader("index", "x", "y", "x+y", "x-y", "y^2-x", "facts(x+y)", "facts(x-y)", "gcd(x+y, x-y)"))
	)

	for x := 0; x < n; x++ {
		T := (x*x + x) >> 1
		for y := 0; y < n; y++ {
			if S := y * y; T == S {
				tbl.AppendRow(NewRow(index, x, y, x+y, x-y, S-x, factorsStr(x+y), factorsStr(x-y), math.GCD(x+y, x-y)))
				index++
			}
		}
	}

	// t.Fatalf("\n%s\n", tbl.String())
}

func TestSort(t *testing.T) {
	tbl := New("Sorted", FltFmtNoExp, 0).SetHeader(NewHeader("index", "letters", "numbers"))
	for i := 0; i < 10; i++ {
		tbl.AppendRow(NewRow(i, string('a'+byte(rand.Intn(26))), rand.Intn(10)))
	}

	tbl.SortOnCol(1)
	// t.Fatalf("\n%s\n", tbl.String())
}
