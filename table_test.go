package table

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"
	"time"
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
		fmt.Println(tbl.Insert(2, NewRow(2, 0.0)).Append(r4).Set(2, 1, 2.2).String())
		fmt.Println(tbl.Types())
		fmt.Println(tbl.InsertCol(0, "Strings", NewCol("zero", "one", "two", "three", "four")).String())
		fmt.Println(tbl.RemoveCol(0))
		fmt.Println(tbl.String())
		fmt.Println(tbl.AppendCol("Times", NewCol(time.Time{}, time.Time{}, time.Time{}, time.Time{}, time.Time{})).String())
	}

	type test struct {
		colNames []string
		ints     []int
		flts     []float64
		bools    []bool
		times    []time.Time
		strs     []string
		exp      *Table
	}

	tests := []test{
		test{
			colNames: []string{"Integers", "Floats", "Booleans", "Times", "Strings"},
			ints:     []int{0, 1, 2, 3, 4},
			flts:     []float64{0, 1, 2, 3, 4},
			bools:    []bool{false, true, false, true, false},
			times:    []time.Time{time.Time{}.AddDate(0, 0, 0), time.Time{}.AddDate(0, 0, 1), time.Time{}.AddDate(0, 0, 2), time.Time{}.AddDate(0, 0, 3), time.Time{}.AddDate(0, 0, 4)},
			strs:     []string{"zero", "one", "two", "three", "four"},
			exp: &Table{
				header: Header{"Integers", "Floats", "Booleans", "Times", "Strings"},
				types:  Types{1, 2, 3, 4, 5},
				body: Body{
					int(0), float64(0.0), false, time.Time{}.AddDate(0, 0, 0), "zero",
					int(1), float64(1.0), true, time.Time{}.AddDate(0, 0, 1), "one",
					int(2), float64(2.0), false, time.Time{}.AddDate(0, 0, 2), "two",
					int(3), float64(3.0), true, time.Time{}.AddDate(0, 0, 3), "three",
					int(4), float64(4.0), false, time.Time{}.AddDate(0, 0, 4), "four",
				},
			},
		},
	}

	for _, tst := range tests {
		rec := New(NewHeader(tst.colNames...))
		for i := 0; i < len(tst.ints); i++ {
			rec.Append(NewRow(tst.ints[i], tst.flts[i], tst.bools[i], tst.times[i], tst.strs[i]))
		}

		if !tst.exp.Equal(rec) {
			t.Errorf("\nexpected %v\nreceived %v\n", tst.exp, rec)
		}
	}
}

func TestImportExport(t *testing.T) {
	inFile, err := os.Open("test0.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer inFile.Close()

	tbl, err := FromCSV(csv.NewReader(inFile))
	if err != nil {
		t.Fatal(err)
	}

	tbl.Insert(1, NewRow(1, 1.1, "one", false))
	tbl.Insert(2, NewRow(2, 2.2, "two", false))
	tbl.Append(NewRow(4, 4.4, "four", false))

	outFile, err := os.OpenFile("test1.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	if err := tbl.ToCSV(csv.NewWriter(outFile)); err != nil {
		t.Fatal(err)
	}

	b, err := tbl.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	t.Errorf("\n%s\n", string(b))
	newTbl := New(nil)
	if err := newTbl.UnmarshalJSON(b); err != nil {
		t.Fatal(err)
	}

	t.Errorf("\n%s\n", newTbl.String())
}

func BenchmarkJSON(b *testing.B) {

}
