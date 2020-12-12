package table

import (
	"encoding/csv"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tidwall/sjson"
)

// TODO: All the testing.

func TestCSV(t *testing.T) {
	var (
		fileName = "test.csv"
		expLines = [][]string{
			[]string{"Integers", "Floats", "Booleans", "Times", "Strings"},
			[]string{"0", "0.0", "false", "0001-01-01T00:00:00Z", "zero"},
			[]string{"1", "1.1", "false", "0001-01-01T00:00:00.000000001Z", "one"},
			[]string{"2", "2.2", "false", "0001-01-01T00:00:00.000000002Z", "two"},
			[]string{"3", "3.3", "true", "0001-01-01T00:00:00.000000003Z", "three"},
			[]string{"4", "4.4", "true", "0001-01-01T00:00:00.000000004Z", "four"},
		}
	)

	{
		// 1. Create a new file with csv package. This file is considered valid before table reads the file's contents.
		os.Remove(fileName)
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		if err := csv.NewWriter(file).WriteAll(expLines[:4]); err != nil {
			t.Fatal(err)
		}

		file.Close()
	}

	{
		// 2. Read from an existing csv file, append several rows, then replace the file's contents.
		file, err := os.OpenFile(fileName, os.O_RDWR, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		tbl, err := FromCSV(csv.NewReader(file))
		if err != nil {
			t.Fatal(err)
		}

		tbl.Append(
			NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
			NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
		)

		file.Close()
		file, err = os.OpenFile(fileName, os.O_RDWR, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		if err := tbl.ToCSV(csv.NewWriter(file)); err != nil {
			t.Fatal(err)
		}

		file.Close()
	}

	{
		// 3. Read table-modified csv file into new table and compare to expected table.
		file, err := os.Open(fileName)
		if err != nil {
			t.Fatal(err)
		}

		tbl, err := FromCSV(csv.NewReader(file))
		if err != nil {
			t.Fatal(err)
		}

		expTbl := New(
			NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
			NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
			NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
			NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
			NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
			NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
		)

		if !expTbl.Equal(tbl) {
			t.Fatalf("\nexpected: %s\nreceived: %s\n", expTbl, tbl)
		}

		file.Close()
	}

	{
		// 4. Compare csv file with csv package.
		file, err := os.Open(fileName)
		if err != nil {
			t.Fatal(err)
		}

		lines, err := csv.NewReader(file).ReadAll()
		if err != nil {
			t.Fatal(err)
		}

		if len(expLines) != len(lines) {
			t.Fatalf("\nexpected %d lines\nreceived %d lines\n", len(expLines), len(lines))
		}

		for i := 0; i < len(expLines); i++ {
			if len(expLines[i]) != len(lines[i]) {
				t.Fatalf("\nexpected line %d to have %d columns\nreceived %d columns\n", i, len(expLines[i]), len(lines[i]))
			}

			for j := 0; j < len(expLines[i]); j++ {
				if !strings.EqualFold(expLines[i][j], lines[i][j]) {
					t.Fatalf("\nexpected line %d, column %d to be %q\nreceived %q\n", i, j, expLines[i][j], lines[i][j])
				}
			}
		}

		file.Close()
	}

	{
		// 5. Clean up.
		if err := os.Remove(fileName); err != nil {
			t.Fatal(err)
		}
	}
}

func TestJSON(t *testing.T) {
	expJSON, err := sjson.Set("", "header", []string{"Integers", "Floats", "Booleans", "Times", "Strings"})
	if err != nil {
		t.Fatal(err)
	}

	expJSON, err = sjson.Set(expJSON, "types", []int{1, 2, 3, 4, 5})
	if err != nil {
		t.Fatal(err)
	}

	expBody := []interface{}{
		0, 0.0, false, "0001-01-01T00:00:00Z", "zero",
		1, 1.1, false, "0001-01-01T00:00:00.000000001Z", "one",
		2, 2.2, false, "0001-01-01T00:00:00.000000002Z", "two",
	}

	expJSON, err = sjson.Set(expJSON, "body", expBody)
	if err != nil {
		t.Fatal(err)
	}

	tbl := FromJSON(expJSON).Append(
		NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
		NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
	)

	expTbl := New(
		NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
		NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
		NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
		NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
		NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
		NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
	)

	if !expTbl.Equal(tbl) {
		t.Fatalf("\nexpected: %s\nreceived: %s\n", expTbl, tbl)
	}

	expBody = append(
		expBody,
		3, 3.3, true, "0001-01-01T00:00:00.000000003Z", "three",
		4, 4.4, true, "0001-01-01T00:00:00.000000004Z", "four",
	)

	expJSON, err = sjson.Set(expJSON, "body", expBody)
	if err != nil {
		t.Fatal(err)
	}

	outJSON := tbl.ToJSON()
	if !strings.EqualFold(expJSON, outJSON) {
		t.Fatalf("\nexpected %q\nreceived %q\n", expJSON, outJSON)
	}
}
