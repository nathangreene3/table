package table

import "testing"

func TestImportCSV(t *testing.T) {
	table, err := ImportCSV("./test0.csv", "star wars")
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Fatalf("%v", table)
}
