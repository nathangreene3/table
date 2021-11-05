package table

var (
	// Fmt0 ...
	//  Int  Str
	// ----------
	//    0  a
	//    1  b
	//    2  c
	Fmt0 = Format{
		UpperHoriz:           "",
		UpperLeftHorizDelim:  "",
		UpperMidHorizDelim:   "",
		UpperRightHorizDelim: "",

		MiddleHoriz:           "-",
		MiddleLeftHorizDelim:  "",
		MiddleMidHorizDelim:   "",
		MiddleRightHorizDelim: "",

		BottomHoriz:           "",
		BottomLeftHorizDelim:  "",
		BottomMidHorizDelim:   "",
		BottomRightHorizDelim: "",

		HeaderLeftDelim:  "",
		HeaderMidDelim:   "",
		HeaderRightDelim: "",

		RowLeftDelim:  "",
		RowMidDelim:   "",
		RowRightDelim: "",
	}

	// Fmt1 ...
	// ----------
	//  Int  Str
	// ----------
	//    0  a
	//    1  b
	//    2  c
	Fmt1 = Format{
		UpperHoriz:           "-",
		UpperLeftHorizDelim:  "",
		UpperMidHorizDelim:   "",
		UpperRightHorizDelim: "",

		MiddleHoriz:           "-",
		MiddleLeftHorizDelim:  "",
		MiddleMidHorizDelim:   "",
		MiddleRightHorizDelim: "",

		BottomHoriz:           "",
		BottomLeftHorizDelim:  "",
		BottomMidHorizDelim:   "",
		BottomRightHorizDelim: "",

		HeaderLeftDelim:  "",
		HeaderMidDelim:   "",
		HeaderRightDelim: "",

		RowLeftDelim:  "",
		RowMidDelim:   "",
		RowRightDelim: "",
	}

	// Fmt2 ...
	//  Int  Str
	// -----------
	//    0  a
	//    1  b
	//    2  c
	// -----------
	Fmt2 = Format{
		UpperHoriz:           "",
		UpperLeftHorizDelim:  "",
		UpperMidHorizDelim:   "",
		UpperRightHorizDelim: "",

		MiddleHoriz:           "-",
		MiddleLeftHorizDelim:  "",
		MiddleMidHorizDelim:   "",
		MiddleRightHorizDelim: "",

		BottomHoriz:           "-",
		BottomLeftHorizDelim:  "",
		BottomMidHorizDelim:   "",
		BottomRightHorizDelim: "",

		HeaderLeftDelim:  "",
		HeaderMidDelim:   "",
		HeaderRightDelim: "",

		RowLeftDelim:  "",
		RowMidDelim:   "",
		RowRightDelim: "",
	}

	// Fmt3 ...
	// ----------
	//  Int  Str
	// ----------
	//    0  a
	//    1  b
	//    2  c
	// ----------
	Fmt3 = Format{
		UpperHoriz:           "-",
		UpperLeftHorizDelim:  "",
		UpperMidHorizDelim:   "",
		UpperRightHorizDelim: "",

		MiddleHoriz:           "-",
		MiddleLeftHorizDelim:  "",
		MiddleMidHorizDelim:   "",
		MiddleRightHorizDelim: "",

		BottomHoriz:           "-",
		BottomLeftHorizDelim:  "",
		BottomMidHorizDelim:   "",
		BottomRightHorizDelim: "",

		HeaderLeftDelim:  "",
		HeaderMidDelim:   "",
		HeaderRightDelim: "",

		RowLeftDelim:  "",
		RowMidDelim:   "",
		RowRightDelim: "",
	}

	// Fmt4 ...
	//  Int   Str
	// ----- -----
	//    0   a
	//    1   b
	//    2   c
	Fmt4 = Format{
		UpperHoriz:           "",
		UpperLeftHorizDelim:  "",
		UpperMidHorizDelim:   "",
		UpperRightHorizDelim: "",

		MiddleHoriz:           "-",
		MiddleLeftHorizDelim:  "",
		MiddleMidHorizDelim:   " ",
		MiddleRightHorizDelim: "",

		BottomHoriz:           "",
		BottomLeftHorizDelim:  "",
		BottomMidHorizDelim:   "",
		BottomRightHorizDelim: "",

		HeaderLeftDelim:  "",
		HeaderMidDelim:   " ",
		HeaderRightDelim: "",

		RowLeftDelim:  "",
		RowMidDelim:   " ",
		RowRightDelim: "",
	}

	// Fmt5 ...
	// +-----+-----+
	// | Int | Str |
	// +-----+-----+
	// |   0 | a   |
	// |   1 | b   |
	// |   2 | c   |
	// +-----+-----+
	Fmt5 = Format{
		UpperHoriz:           "-",
		UpperLeftHorizDelim:  "+",
		UpperMidHorizDelim:   "+",
		UpperRightHorizDelim: "+",

		MiddleHoriz:           "-",
		MiddleLeftHorizDelim:  "+",
		MiddleMidHorizDelim:   "+",
		MiddleRightHorizDelim: "+",

		BottomHoriz:           "-",
		BottomLeftHorizDelim:  "+",
		BottomMidHorizDelim:   "+",
		BottomRightHorizDelim: "+",

		HeaderLeftDelim:  "|",
		HeaderMidDelim:   "|",
		HeaderRightDelim: "|",

		RowLeftDelim:  "|",
		RowMidDelim:   "|",
		RowRightDelim: "|",
	}
)

// Format holds decoration characters for displaying a table.
type Format struct {
	UpperHoriz           string
	UpperLeftHorizDelim  string
	UpperMidHorizDelim   string
	UpperRightHorizDelim string

	MiddleHoriz           string
	MiddleLeftHorizDelim  string
	MiddleMidHorizDelim   string
	MiddleRightHorizDelim string

	BottomHoriz           string
	BottomLeftHorizDelim  string
	BottomMidHorizDelim   string
	BottomRightHorizDelim string

	HeaderLeftDelim  string
	HeaderMidDelim   string
	HeaderRightDelim string

	RowLeftDelim  string
	RowMidDelim   string
	RowRightDelim string
}
