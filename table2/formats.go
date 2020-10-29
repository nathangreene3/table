package table2

// Formats ...
type Formats []string

// Equal ...
func (f Formats) Equal(fmts Formats) bool {
	if len(f) != len(fmts) {
		return false
	}

	var i int
	for ; i < len(f) && f[i] == fmts[i]; i++ {
	}

	return i == len(f)
}
