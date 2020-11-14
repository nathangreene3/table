package table2

import (
	"encoding/csv"
	"fmt"
	"os"
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
		fmt.Println(tbl.Insert(2, NewRow(2, 0.0)).Append(r4).Set(2, 1, 2.2).String())
		fmt.Println(tbl.Fmts())
		fmt.Println(tbl.InsertCol(0, "Strings", NewCol("zero", "one", "two", "three", "four")).String())
		fmt.Println(tbl.RemoveCol(0))
		fmt.Println(tbl.String())
	}

	t.Fatal()
}

func TestImportExport(t *testing.T) {
	inFile, err := os.Open("test0.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer inFile.Close()

	tbl, err := Import(csv.NewReader(inFile))
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

	if err := tbl.Export(csv.NewWriter(outFile)); err != nil {
		t.Fatal(err)
	}

	b, err := tbl.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	t.Errorf("\n%s\n", string(b))
}

func BenchmarkJSON(b *testing.B) {

}
