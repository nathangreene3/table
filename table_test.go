package table

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/tidwall/sjson"
)

// TODO: All the testing.

func TestFormat(t *testing.T) {
	tests := []struct {
		tbl *Table
		fmt Format
		exp string
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
			tbl: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
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
			tbl: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
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
			tbl: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
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
			tbl: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
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
			tbl: New(
				NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
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
			t.Errorf("\n"+
				"expected:\n%q\n"+
				"received:\n%q\n",
				test.exp,
				rec,
			)
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
			tbl: New(NewHeader()),
			exp: New(
				NewHeader("Integers"),
				NewRow(0),
				NewRow(1),
				NewRow(2),
				NewRow(3),
				NewRow(4),
			),
			colName: "Integers",
			col:     NewCol(0, 1, 2, 3, 4),
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
			t.Errorf("\n"+
				"expected:\n%v\n"+
				"received:\n%v\n",
				test.exp,
				rec,
			)
		}
	}
}

func TestCSV(t *testing.T) {
	const fileName = "test.csv"
	var expLines = [][]string{
		{"Integers", "Floats", "Booleans", "Times", "Strings"},
		{"0", "0.0", "false", "0001-01-01T00:00:00Z", "zero"},
		{"1", "1.1", "false", "0001-01-01T00:00:00.000000001Z", "one"},
		{"2", "2.2", "false", "0001-01-01T00:00:00.000000002Z", "two"},
		{"3", "3.3", "true", "0001-01-01T00:00:00.000000003Z", "three"},
		{"4", "4.4", "true", "0001-01-01T00:00:00.000000004Z", "four"},
	}

	// 1. Create a new file with csv package. This file is considered valid before table reads the file's contents.
	os.Remove(fileName) // Intentionally ignore error here

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	var closed, removed bool
	defer func() {
		if !closed {
			if err := file.Close(); err != nil {
				t.Fatal(err)
			}

			closed = true
		}

		if !removed {
			if err := os.Remove(fileName); err != nil {
				t.Fatal(err)
			}

			removed = true
		}
	}()

	if err := csv.NewWriter(file).WriteAll(expLines[:4]); err != nil {
		t.Fatal(err)
	}

	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	closed = true

	// 2. Read from an existing csv file, append several rows, then replace the file's contents.
	tbl, err := FromCSV(fileName)
	if err != nil {
		t.Fatal(err)
	}

	tbl.Append(
		NewRow(3, 3.3, true, NewFTime(time.Time{}.Add(3)), "three"),
		NewRow(4, 4.4, true, NewFTime(time.Time{}.Add(4)), "four"),
	)

	if err := tbl.WriteCSV(fileName); err != nil {
		t.Fatal(err)
	}

	// 3. Read table-modified csv file into new table and compare to expected table.
	tbl, err = FromCSV(fileName)
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
		t.Fatalf("\n"+
			"expected: %s\n"+
			"received: %s\n",
			expTbl,
			tbl,
		)
	}

	// 4. Compare csv file with csv package.
	file, err = os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}

	closed = false

	defer func() {
		if !closed {
			if err := file.Close(); err != nil {
				t.Fatal(err)
			}

			closed = true
		}

		if !removed {
			if err := os.Remove(fileName); err != nil {
				t.Fatal(err)
			}

			removed = true
		}
	}()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(expLines) != len(lines) {
		t.Fatalf("\n"+
			"expected %d lines\n"+
			"received %d lines\n",
			len(expLines),
			len(lines),
		)
	}

	for i := 0; i < len(expLines); i++ {
		if len(expLines[i]) != len(lines[i]) {
			t.Fatalf("\n"+
				"expected line %d to have %d columns\n"+
				"received %d columns\n",
				i, len(expLines[i]),
				len(lines[i]),
			)
		}

		for j := 0; j < len(expLines[i]); j++ {
			if !strings.EqualFold(expLines[i][j], lines[i][j]) {
				t.Fatalf("\n"+
					"expected line %d, column %d to be %q\n"+
					"received %q\n",
					i, j, expLines[i][j],
					lines[i][j],
				)
			}
		}
	}

	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	closed = true
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
			t.Errorf("\n"+
				"expected (%d,%d)\n"+
				"received (%d,%d)\n",
				test.expM, test.expN,
				recM, recN,
			)
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
			rec              = Generate(hdr, 5, recGen).Filter(fltr)
			exp              = Generate(hdr, 3, expGen)
		)

		if !exp.Equal(rec) {
			t.Fatalf("\n"+
				"expected %s\n"+
				"received %s\n",
				exp,
				rec,
			)
		}
	}

	{
		// Primes
		// var (
		// 	hdr            = NewHeader("n")
		// 	gen  generator = func(i int) Row { return NewRow(i + 1) }
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

		// t.Fatal(gen(hdr, 100, gen).Filter(fltr))
	}
}

func TestGen(t *testing.T) {
	intStr := func(n int) string {
		var s string
		switch n {
		case 0:
			s = "zero"
		case 1:
			s = "one"
		case 2:
			s = "two"
		case 3:
			s = "three"
		case 4:
			s = "four"
		case 5:
			s = "five"
		case 6:
			s = "six"
		case 7:
			s = "seven"
		case 8:
			s = "eight"
		case 9:
			s = "nine"
		default:
			panic("Value not yet supported")
		}

		return s
	}

	tests := []struct {
		h   Header
		m   int
		f   Generator
		exp *Table
	}{
		{
			h: NewHeader("Integers", "Floats", "Booleans", "Times", "Strings"),
			m: 5,
			f: func(i int) Row {
				s := strconv.Itoa(i)
				f, err := strconv.ParseFloat(s+"."+s, 64)
				if err != nil {
					panic(err)
				}

				return NewRow(i, f, 2 < i, NewFTime(time.Time{}.Add(time.Duration(i))), intStr(i))
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

	for _, test := range tests {
		if rec := Generate(test.h, test.m, test.f); !test.exp.Equal(rec) {
			t.Errorf("\n"+
				"expected:\n%v\n"+
				"received:\n%v\n",
				test.exp,
				rec,
			)
		}
	}
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
			t.Errorf("\n"+
				" test %d: Join\n"+
				"expected: \n%s\n"+
				"received: \n%s\n",
				i,
				test.exp,
				rec)
		}

		if rec := join1(test.tbl...); !test.exp.Equal(rec) {
			t.Errorf("\nTest %d: join\nexpected: \n%s\nreceived: \n%s\n", i, test.exp, rec)
		}

		if rec := join2(test.tbl...); !test.exp.Equal(rec) {
			t.Errorf("\n"+
				" test %d: join2\n"+
				"expected: \n%s\n"+
				"received: \n%s\n",
				i,
				test.exp,
				rec,
			)
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
		t.Fatalf("\n"+
			"expected: %s\n"+
			"received: %s\n",
			expTbl,
			tbl,
		)
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

	outJSON := tbl.JSON()
	if !strings.EqualFold(expJSON, outJSON) {
		t.Fatalf("\n"+
			"expected %q\n"+
			"received %q\n",
			expJSON,
			outJSON,
		)
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
			t.Errorf("\n"+
				"expected %q\n"+
				"received %q\n",
				string(expBts),
				string(recBts),
			)
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
			rec              = Generate(hdr, 5, recGen).Map(mpr)
			exp              = Generate(hdr, 5, expGen)
		)

		if !exp.Equal(rec) {
			t.Fatalf("\n"+
				"expected %s\n"+
				"received %s\n",
				exp,
				rec,
			)
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
			rec              = Generate(h, n, recGen).Reduce(rdcr)
			exp              = NewRow(n * (n - 1) >> 1)
		)

		if !exp.Equal(rec) {
			t.Fatalf("\n"+
				"expected %s\n"+
				"received %s\n",
				exp,
				rec,
			)
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
			t.Errorf("\n"+
				"   given: col = %d\n"+
				"expected:\n%v\n"+
				"received:\n%v\n",
				test.col,
				test.exp,
				rec,
			)
		}
	}
}

// ------------------------------------------------------------------------------------
// Benchmarks
// ------------------------------------------------------------------------------------

func BenchmarkJoin(b *testing.B) {
	// Linear benchmarks
	var (
		m0, m1, dm = 8, 8, 1 // Number of rows
		n0, n1, dn = 8, 8, 1 // Number of columns
		k0, k1, dk = 8, 8, 1 // Number of tables
		s          = strings.Repeat("hello", 8)
	)

	for n := n0; n <= n1; n += dn {
		for m := m0; m <= m1; m += dm {
			col := NewCol()
			for mi := 0; mi < m; mi++ {
				col = append(col, s)
			}

			tbl := New(NewHeader())
			for ni := 1; ni < n; ni++ {
				tbl.AppendCol("Strings", col)
			}

			tbls := make([]*Table, 0, k1)
			for k := k0; k <= k1; k += dk {
				for ; len(tbls) < k; tbls = append(tbls, tbl.Copy()) {
				}

				benchmarkJoin(b, fmt.Sprintf("Join %d %dx%d tables", k, m, n), Join, tbls...)
			}
		}
	}

	for n := n0; n <= n1; n += dn {
		for m := m0; m <= m1; m += dm {
			col := NewCol()
			for mi := 0; mi < m; mi++ {
				col = append(col, s)
			}

			tbl := New(NewHeader())
			for ni := 1; ni < n; ni++ {
				tbl.AppendCol("Strings", col)
			}

			tbls := make([]*Table, 0, k1)
			for k := k0; k <= k1; k += dk {
				for ; len(tbls) < k; tbls = append(tbls, tbl.Copy()) {
				}

				benchmarkJoin(b, fmt.Sprintf("join1 %d %dx%d tables", k, m, n), join1, tbls...)
			}
		}
	}

	for n := n0; n <= n1; n += dn {
		for m := m0; m <= m1; m += dm {
			col := NewCol()
			for mi := 0; mi < m; mi++ {
				col = append(col, s)
			}

			tbl := New(NewHeader())
			for ni := 1; ni < n; ni++ {
				tbl.AppendCol("Strings", col)
			}

			tbls := make([]*Table, 0, k1)
			for k := k0; k <= k1; k += dk {
				for ; len(tbls) < k; tbls = append(tbls, tbl.Copy()) {
				}

				benchmarkJoin(b, fmt.Sprintf("join2 %d %dx%d tables", k, m, n), join2, tbls...)
			}
		}
	}
}

func benchmarkJoin(b *testing.B, name string, joinFn func(tbl ...*Table) *Table, tbls ...*Table) bool {
	g := func(b *testing.B) {
		var tbl *Table
		for i := 0; i < b.N; i++ {
			tbl = joinFn(tbls...)
		}
		_ = tbl
	}

	return b.Run(name, g)
}
