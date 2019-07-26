package table

import "testing"

func TestImportCSV(t *testing.T) {
	table, err := ImportCSV("./test0.csv", "star wars", 'f', 3)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	table.Clean()
	if err = table.ExportCSV("test.csv"); err != nil {
		t.Fatalf("\n%v", err)
	}

	t.Fatalf("\n%s", table.String())
}
