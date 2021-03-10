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

func TestFormat(t *testing.T) {
	tests := []struct {
		tbl *Table
		fmt format
		exp string
	}{
		{
			tbl: New(NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			fmt: Fmt0,
			exp: "\n" +
				" Integers  Floats  Booleans  Times                           Strings \n" +
				"---------------------------------------------------------------------\n" +
				"        0     0.0  false     0001-01-01T00:00:00Z            zero    \n" +
				"        1     1.1  false     0001-01-01T00:00:00.000000001Z  one     \n" +
				"        2     2.2  false     0001-01-01T00:00:00.000000002Z  two     \n" +
				"        3     3.3  true      0001-01-01T00:00:00.000000003Z  three   \n" +
				"        4     4.4  true      0001-01-01T00:00:00.000000004Z  four    \n",
		},
		{
			tbl: New(NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			fmt: Fmt1,
			exp: "\n" +
				"---------------------------------------------------------------------\n" +
				" Integers  Floats  Booleans  Times                           Strings \n" +
				"---------------------------------------------------------------------\n" +
				"        0     0.0  false     0001-01-01T00:00:00Z            zero    \n" +
				"        1     1.1  false     0001-01-01T00:00:00.000000001Z  one     \n" +
				"        2     2.2  false     0001-01-01T00:00:00.000000002Z  two     \n" +
				"        3     3.3  true      0001-01-01T00:00:00.000000003Z  three   \n" +
				"        4     4.4  true      0001-01-01T00:00:00.000000004Z  four    \n",
		},
		{
			tbl: New(NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			fmt: Fmt2,
			exp: "\n" +
				" Integers  Floats  Booleans  Times                           Strings \n" +
				"---------------------------------------------------------------------\n" +
				"        0     0.0  false     0001-01-01T00:00:00Z            zero    \n" +
				"        1     1.1  false     0001-01-01T00:00:00.000000001Z  one     \n" +
				"        2     2.2  false     0001-01-01T00:00:00.000000002Z  two     \n" +
				"        3     3.3  true      0001-01-01T00:00:00.000000003Z  three   \n" +
				"        4     4.4  true      0001-01-01T00:00:00.000000004Z  four    \n" +
				"---------------------------------------------------------------------\n",
		},
		{
			tbl: New(NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			fmt: Fmt3,
			exp: "\n" +
				"---------------------------------------------------------------------\n" +
				" Integers  Floats  Booleans  Times                           Strings \n" +
				"---------------------------------------------------------------------\n" +
				"        0     0.0  false     0001-01-01T00:00:00Z            zero    \n" +
				"        1     1.1  false     0001-01-01T00:00:00.000000001Z  one     \n" +
				"        2     2.2  false     0001-01-01T00:00:00.000000002Z  two     \n" +
				"        3     3.3  true      0001-01-01T00:00:00.000000003Z  three   \n" +
				"        4     4.4  true      0001-01-01T00:00:00.000000004Z  four    \n" +
				"---------------------------------------------------------------------\n",
		},
		{
			tbl: New(NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			fmt: Fmt4,
			exp: "\n" +
				" Integers   Floats   Booleans   Times                            Strings \n" +
				"---------- -------- ---------- -------------------------------- ---------\n" +
				"        0      0.0   false      0001-01-01T00:00:00Z             zero    \n" +
				"        1      1.1   false      0001-01-01T00:00:00.000000001Z   one     \n" +
				"        2      2.2   false      0001-01-01T00:00:00.000000002Z   two     \n" +
				"        3      3.3   true       0001-01-01T00:00:00.000000003Z   three   \n" +
				"        4      4.4   true       0001-01-01T00:00:00.000000004Z   four    \n",
		},
		{
			tbl: New(NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			fmt: Fmt5,
			exp: "\n" +
				"+----------+--------+----------+--------------------------------+---------+\n" +
				"| Integers | Floats | Booleans | Times                          | Strings |\n" +
				"+----------+--------+----------+--------------------------------+---------+\n" +
				"|        0 |    0.0 | false    | 0001-01-01T00:00:00Z           | zero    |\n" +
				"|        1 |    1.1 | false    | 0001-01-01T00:00:00.000000001Z | one     |\n" +
				"|        2 |    2.2 | false    | 0001-01-01T00:00:00.000000002Z | two     |\n" +
				"|        3 |    3.3 | true     | 0001-01-01T00:00:00.000000003Z | three   |\n" +
				"|        4 |    4.4 | true     | 0001-01-01T00:00:00.000000004Z | four    |\n" +
				"+----------+--------+----------+--------------------------------+---------+\n",
		},
	}

	for _, test := range tests {
		if rec := test.tbl.Format(test.fmt); test.exp != rec {
			t.Errorf("\nexpected:\n%q\nreceived:\n%q\n", test.exp, rec)
		}
	}
}

func TestAppendCol(t *testing.T) {
	tests := []struct {
		tbl, exp *Table
		colName  string
		col      Column
	}{
		{
			tbl:     New(NewHeader()),
			exp:     New(NewHeader("")),
			colName: "",
			col:     NewCol(),
		},
		{
			tbl: New(
				NewHeader("Integers"),
				NewRow(0),
				NewRow(1),
			),
			exp: New(
				NewHeader("Integers", "Floats"),
				NewRow(0, 0.0),
				NewRow(1, 1.1),
			),
			colName: "Floats",
			col:     NewCol(0.0, 1.1),
		},
	}

	for _, test := range tests {
		if rec := test.tbl.AppendCol(test.colName, test.col); !test.exp.Equal(rec) {
			t.Errorf("\nexpected:\n%v\nreceived:\n%v\n", test.exp, rec)
		}
	}
}

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
		os.Remove(fileName) // Intentionally ignore error here
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
		tbl, err := FromCSV(fileName)
		if err != nil {
			t.Fatal(err)
		}

		tbl.Append(
			NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
			NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
		)

		if err := tbl.ToCSV(fileName); err != nil {
			t.Fatal(err)
		}
	}

	{
		// 3. Read table-modified csv file into new table and compare to expected table.
		tbl, err := FromCSV(fileName)
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

func TestDims(t *testing.T) {
	tests := []struct {
		tbl        *Table
		expM, expN int
	}{
		{
			tbl: New(NewHeader()),
		},
		{
			tbl: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			expM: 5,
			expN: 5,
		},
	}

	for _, test := range tests {
		if recM, recN := test.tbl.Dims(); test.expM != recM || test.expN != recN {
			t.Errorf("\nexpected (%d,%d)\nreceived (%d,%d)\n", test.expM, test.expN, recM, recN)
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

	{
		// Primes
		// var (
		// 	hdr            = NewHeader("n")
		// 	gen  Generator = func(i int) Row { return NewRow(i + 1) }
		// 	fltr Filterer  = func(r Row) bool {
		// 		n := r[0].(int)
		// 		if n == 2 {
		// 			return true
		// 		}

		// 		if n&1 == 0 {
		// 			return false
		// 		}

		// 		for i := 3; i < n; i += 2 {
		// 			if n%i == 0 {
		// 				return false
		// 			}
		// 		}

		// 		return true
		// 	}
		// )

		// t.Fatal(Gen(hdr, 100, gen).Filter(fltr))
	}
}

func TestGen(t *testing.T) {
	// TODO
}

func TestJoin(t *testing.T) {
	tests := []struct {
		tbl []*Table
		exp *Table
	}{
		{
			tbl: []*Table{New(NewHeader()), New(NewHeader()), New(NewHeader())},
			exp: New(NewHeader()),
		},
		{
			tbl: []*Table{
				New(
					NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
					NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
					NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
					NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
					NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
					NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
				),
			},
			exp: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
		},
		{
			tbl: []*Table{
				New(
					NewHeader("Integers"),
					NewRow(0),
					NewRow(1),
					NewRow(2),
					NewRow(3),
					NewRow(4),
				),
				New(
					NewHeader("Floats", "Booleans", "Times", "Strings"),
					NewRow(0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
					NewRow(1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
					NewRow(2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
					NewRow(3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
					NewRow(4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
				),
			},
			exp: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
		},
		{
			tbl: []*Table{
				New(
					NewHeader("Integers", "Floats", "Booleans", "Times"),
					NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0))),
					NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1))),
					NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2))),
					NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3))),
					NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4))),
				),
				New(
					NewHeader("Strings"),
					NewRow("zero"),
					NewRow("one"),
					NewRow("two"),
					NewRow("three"),
					NewRow("four"),
				),
			},
			exp: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
		},
		{
			tbl: []*Table{
				New(
					NewHeader("Integers"),
					NewRow(0),
					NewRow(1),
					NewRow(2),
					NewRow(3),
					NewRow(4),
				),
				New(
					NewHeader("Floats", "Booleans", "Times"),
					NewRow(0.0, false, NewFTime(time.Time{}.Add(0))),
					NewRow(1.1, false, NewFTime(time.Time{}.Add(1))),
					NewRow(2.2, false, NewFTime(time.Time{}.Add(2))),
					NewRow(3.3, true, NewFTime(time.Time{}.Add(3))),
					NewRow(4.4, true, NewFTime(time.Time{}.Add(4))),
				),
				New(
					NewHeader("Strings"),
					NewRow("zero"),
					NewRow("one"),
					NewRow("two"),
					NewRow("three"),
					NewRow("four"),
				),
			},
			exp: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
		},
	}

	for i, test := range tests {
		if rec := Join(test.tbl...); !test.exp.Equal(rec) {
			t.Errorf("\nTest %d: Join\nexpected: \n%s\nreceived: \n%s\n", i, test.exp, rec)
		}

		if rec := join(test.tbl...); !test.exp.Equal(rec) {
			t.Errorf("\nTest %d: join\nexpected: \n%s\nreceived: \n%s\n", i, test.exp, rec)
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
		// Sum 0 + 1 + ... + (n-1) = n*(n-1)/2
		var (
			h                = NewHeader("n")
			recGen Generator = func(i int) Row { return Row{i} }
			rdcr   Reducer   = func(dst, src Row) { dst[0] = dst[0].(int) + src[0].(int) }
			n                = 5
			rec              = Gen(h, n, recGen).Reduce(rdcr)
			exp              = NewRow(n * (n - 1) >> 1)
		)

		if !exp.Equal(rec) {
			t.Fatalf("\nexpected %s\nreceived %s\n", exp, rec)
		}
	}
}

func TestStable(t *testing.T) {
	tests := []struct {
		tbl, exp *Table
		col      int
	}{
		{
			tbl: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			exp: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			col: 0,
		},
		{
			tbl: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			exp: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			col: 1,
		},
		{
			tbl: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			exp: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			col: 2,
		},
		{
			tbl: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			exp: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			col: 3,
		},
		{
			tbl: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
			),
			exp: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
				NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
				NewRow(1, 1.1, false, NewFTime(time.Time{}.Add(1)), "one"),
				NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
				NewRow(2, 2.2, false, NewFTime(time.Time{}.Add(2)), "two"),
				NewRow(0, 0.0, false, NewFTime(time.Time{}.Add(0)), "zero"),
			),
			col: 4,
		},
	}

	for _, test := range tests {
		if rec := test.tbl.Copy().Stable(test.col); !test.exp.Equal(rec) {
			t.Errorf("\n   given: col = %d\nexpected:\n%v\nreceived:\n%v\n", test.col, test.exp, rec)
		}
	}
}
