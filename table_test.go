package table

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
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
			{"Integers", "Floats", "Booleans", "Times", "Strings"},
			{"0", "0.0", "false", "0001-01-01T00:00:00Z", "zero"},
			{"1", "1.1", "false", "0001-01-01T00:00:00.000000001Z", "one"},
			{"2", "2.2", "false", "0001-01-01T00:00:00.000000002Z", "two"},
			{"3", "3.3", "true", "0001-01-01T00:00:00.000000003Z", "three"},
			{"4", "4.4", "true", "0001-01-01T00:00:00.000000004Z", "four"},
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

func TestFilter(t *testing.T) {
	{
		// Evens
		var (
			hdr              = NewHeader("n", "even")
			recGen Generator = func(i int) Row { return Row{i, i%2 == 0} }
			expGen Generator = func(i int) Row { i <<= 1; return Row{i, i%2 == 0} }
			fltr   Filterer  = func(r Row) bool { return r[0].(int)%2 == 0 }
			rec              = Gen(hdr, 5, recGen).Filter(fltr)
			exp              = Gen(hdr, 3, expGen)
		)

		if !exp.Equal(rec) {
			t.Fatalf("\nexpected %s\nreceived %s\n", exp, rec)
		}
	}
}

func TestJoin(t *testing.T) {
	// TODO
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

	tbl, err := FromJSON(expJSON)
	if err != nil {
		t.Fatal(err)
	}

	tbl.Append(
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

	{
		// Default json interface implementations
		expBts, err := json.Marshal(expTbl)
		if err != nil {
			t.Fatal(err)
		}

		tbl := New(nil)
		if err := json.Unmarshal(expBts, tbl); err != nil {
			t.Fatal(err)
		}

		if !tbl.Equal(expTbl) {
			t.Fatal()
		}

		recBts, err := json.Marshal(tbl)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(expBts, recBts) {
			t.Errorf("\nexpected %q\nreceived %q\n", string(expBts), string(recBts))
		}
	}
}

func TestMap(t *testing.T) {
	{
		// Evens
		var (
			hdr              = NewHeader("n", "is even")
			recGen Generator = func(i int) Row { return Row{i, i%2 == 0} }
			expGen Generator = func(i int) Row { i <<= 1; return Row{i, i%2 == 0} }
			mpr    Mapper    = func(r Row) { r[0] = r[0].(int) << 1; r[1] = r[0].(int)%2 == 0 }
			rec              = Gen(hdr, 5, recGen).Map(mpr)
			exp              = Gen(hdr, 5, expGen)
		)

		if !exp.Equal(rec) {
			t.Fatalf("\nexpected %s\nreceived %s\n", exp, rec)
		}
	}
}

func TestReduce(t *testing.T) {
	{
		// Sum
		var (
			hdr              = NewHeader("n")
			recGen Generator = func(i int) Row { return Row{i} }
			rdcr   Reducer   = func(dst, src Row) { dst[0] = dst[0].(int) + src[0].(int) }
			rec              = Gen(hdr, 5, recGen).Reduce(rdcr)
			exp              = NewRow(10)
		)

		if !exp.Equal(rec) {
			t.Fatalf("\nexpected %s\nreceived %s\n", exp, rec)
		}
	}
}
