package table

import "testing"

func TestImportCSV(t *testing.T) {
	table, err := ImportCSV("./test0.csv", "star wars", 'f', 3)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if err = table.ExportCSV("test1.csv"); err != nil {
		t.Fatalf("\n%v", err)
	}

	t.Fatalf("\n%s", table.String())
}
