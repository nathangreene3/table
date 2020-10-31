package table2

const (
	// Bool ...
	Bool Format = "bool"

	// Flt ...
	Flt Format = "float64"

	// Int ...
	Int Format = "int"

	// Str ...
	Str Format = "string"
)

// Format ...
type Format string

// Formats ...
type Formats []Format

// NewFormats ...
func NewFormats(fmts ...Format) Formats {
	return append(make(Formats, 0, len(fmts)), fmts...)
}

// Copy ...
func (f Formats) Copy() Formats {
	return append(make(Formats, 0, len(f)), f...)
}

// Equal ...
func (f Formats) Equal(fmts Formats) bool {
	if len(f) != len(fmts) {
		return false
	}

	for i := 0; i < len(f); i++ {
		if f[i] != fmts[i] {
			return false
		}
	}

	return true
}

// Fmt ...
func Fmt(x interface{}) Format {
	switch x.(type) {
	case bool:
		return Bool
	case float64:
		return Flt
	case int:
		return Int
	case string:
		return Str
	default:
		return ""
	}
}
