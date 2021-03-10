package table

var (
	// Fmt0 ...
	//  Int  Str
	// ----------
	//    0  a
	//    1  b
	//    2  c
	Fmt0 = format{
		upperHoriz:           "",
		upperLeftHorizDelim:  "",
		upperMidHorizDelim:   "",
		upperRightHorizDelim: "",

		middleHoriz:           "-",
		middleLeftHorizDelim:  "",
		middleMidHorizDelim:   "",
		middleRightHorizDelim: "",

		bottomHoriz:           "",
		bottomLeftHorizDelim:  "",
		bottomMidHorizDelim:   "",
		bottomRightHorizDelim: "",

		headerLeftDelim:  "",
		headerMidDelim:   "",
		headerRightDelim: "",

		rowLeftDelim:  "",
		rowMidDelim:   "",
		rowRightDelim: "",
	}

	// Fmt1 ...
	// ----------
	//  Int  Str
	// ----------
	//    0  a
	//    1  b
	//    2  c
	Fmt1 = format{
		upperHoriz:           "-",
		upperLeftHorizDelim:  "",
		upperMidHorizDelim:   "",
		upperRightHorizDelim: "",

		middleHoriz:           "-",
		middleLeftHorizDelim:  "",
		middleMidHorizDelim:   "",
		middleRightHorizDelim: "",

		bottomHoriz:           "",
		bottomLeftHorizDelim:  "",
		bottomMidHorizDelim:   "",
		bottomRightHorizDelim: "",

		headerLeftDelim:  "",
		headerMidDelim:   "",
		headerRightDelim: "",

		rowLeftDelim:  "",
		rowMidDelim:   "",
		rowRightDelim: "",
	}

	// Fmt2 ...
	//  Int  Str
	// -----------
	//    0  a
	//    1  b
	//    2  c
	// -----------
	Fmt2 = format{
		upperHoriz:           "",
		upperLeftHorizDelim:  "",
		upperMidHorizDelim:   "",
		upperRightHorizDelim: "",

		middleHoriz:           "-",
		middleLeftHorizDelim:  "",
		middleMidHorizDelim:   "",
		middleRightHorizDelim: "",

		bottomHoriz:           "-",
		bottomLeftHorizDelim:  "",
		bottomMidHorizDelim:   "",
		bottomRightHorizDelim: "",

		headerLeftDelim:  "",
		headerMidDelim:   "",
		headerRightDelim: "",

		rowLeftDelim:  "",
		rowMidDelim:   "",
		rowRightDelim: "",
	}

	// Fmt3 ...
	// ----------
	//  Int  Str
	// ----------
	//    0  a
	//    1  b
	//    2  c
	// ----------
	Fmt3 = format{
		upperHoriz:           "-",
		upperLeftHorizDelim:  "",
		upperMidHorizDelim:   "",
		upperRightHorizDelim: "",

		middleHoriz:           "-",
		middleLeftHorizDelim:  "",
		middleMidHorizDelim:   "",
		middleRightHorizDelim: "",

		bottomHoriz:           "-",
		bottomLeftHorizDelim:  "",
		bottomMidHorizDelim:   "",
		bottomRightHorizDelim: "",

		headerLeftDelim:  "",
		headerMidDelim:   "",
		headerRightDelim: "",

		rowLeftDelim:  "",
		rowMidDelim:   "",
		rowRightDelim: "",
	}

	// Fmt4 ...
	//  Int   Str
	// ----- -----
	//    0   a
	//    1   b
	//    2   c
	Fmt4 = format{
		upperHoriz:           "",
		upperLeftHorizDelim:  "",
		upperMidHorizDelim:   "",
		upperRightHorizDelim: "",

		middleHoriz:           "-",
		middleLeftHorizDelim:  "",
		middleMidHorizDelim:   " ",
		middleRightHorizDelim: "",

		bottomHoriz:           "",
		bottomLeftHorizDelim:  "",
		bottomMidHorizDelim:   "",
		bottomRightHorizDelim: "",

		headerLeftDelim:  "",
		headerMidDelim:   " ",
		headerRightDelim: "",

		rowLeftDelim:  "",
		rowMidDelim:   " ",
		rowRightDelim: "",
	}

	// Fmt5 ...
	// +-----+-----+
	// | Int | Str |
	// +-----+-----+
	// |   0 | a   |
	// |   1 | b   |
	// |   2 | c   |
	// +-----+-----+
	Fmt5 = format{
		upperHoriz:           "-",
		upperLeftHorizDelim:  "+",
		upperMidHorizDelim:   "+",
		upperRightHorizDelim: "+",

		middleHoriz:           "-",
		middleLeftHorizDelim:  "+",
		middleMidHorizDelim:   "+",
		middleRightHorizDelim: "+",

		bottomHoriz:           "-",
		bottomLeftHorizDelim:  "+",
		bottomMidHorizDelim:   "+",
		bottomRightHorizDelim: "+",

		headerLeftDelim:  "|",
		headerMidDelim:   "|",
		headerRightDelim: "|",

		rowLeftDelim:  "|",
		rowMidDelim:   "|",
		rowRightDelim: "|",
	}
)

// format holds decoration characters for displaying a table.
type format struct {
	upperHoriz           string
	upperLeftHorizDelim  string
	upperMidHorizDelim   string
	upperRightHorizDelim string

	middleHoriz           string
	middleLeftHorizDelim  string
	middleMidHorizDelim   string
	middleRightHorizDelim string

	bottomHoriz           string
	bottomLeftHorizDelim  string
	bottomMidHorizDelim   string
	bottomRightHorizDelim string

	headerLeftDelim  string
	headerMidDelim   string
	headerRightDelim string

	rowLeftDelim  string
	rowMidDelim   string
	rowRightDelim string
}
