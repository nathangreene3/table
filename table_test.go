package table

import (
	"encoding/csv"
	"os"
	"testing"
)

func TestImportCSV(t *testing.T) {
	inFile, err := os.Open("test0.csv")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer inFile.Close()

	table, err := ImportFromCSV(csv.NewReader(inFile), "star wars", 'f', 3)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	outFile, err := os.OpenFile("test1.csv", os.O_WRONLY, os.ModeAppend)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer outFile.Close()

	if err = table.ExportToCSV(csv.NewWriter(outFile)); err != nil {
		t.Fatalf("\n%v", err)
	}

	// t.Fatalf("\n%s", table.String())
}
