package table

import (
	"encoding/csv"
	"os"
	"testing"
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
