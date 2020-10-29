package table

import (
	"bytes"
	"strings"
)

// toLines converts comma and newline (cleans return carriage) delimited bytes to
// a 2D string slice. For example, []byte("hello,world\r\ngoodbye,cruel world")
// maps to
// [][]string{
// 	[]string{"hello", "world"},
//  []string{"goodbye", "cruel world"}
// }.
func toLines(b []byte) [][]string {
	split := bytes.Split(bytes.ReplaceAll(b, []byte{'\r', '\n'}, []byte{'\n'}), []byte{'\n'})
	lines := make([][]string, 0, len(split))
	for i := 0; i < len(split); i++ {
		lines = append(lines, strings.Split(string(split[i]), ","))
	}

	return lines
}
