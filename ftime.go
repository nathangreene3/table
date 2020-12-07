package table

import "time"

// timeFmts is a list of time formats supported in parsing a time stamp.
var timeFmts = []string{
	// These are ordered in what is thought to be most likely to least likely to be used as a time format
	time.RFC3339Nano, // Default; This MUST precede RFC3339
	time.RFC3339,
	"2006-01-02 15:04:05.999999999 -0700 MST", // time's default format returned in time.String
	time.UnixDate,
	time.ANSIC,
	time.RFC1123Z,
	time.RFC1123,
	time.RFC850,
	time.RFC822Z,
	time.RFC822,
	time.StampNano,
	time.StampMicro,
	time.StampMilli,
	time.Stamp,
	time.Kitchen,
	time.RubyDate,
}

// FTime is a time.Time together with a format that was either
// used in parsing the timestamp from a string or set to
// format the timestamp.
type FTime struct {
	time   time.Time
	format string
}

// NewFTime returns a new FTime. A single format may be passed
// indicating how the underlying time.Time will be formatted
// when writing to string. This panics if more than one format
// is provided. If no format is provided, time.RFC3339Nano is
// set as the format by default.
func NewFTime(t time.Time, format ...string) FTime {
	switch len(format) {
	case 0:
		return FTime{time: t, format: time.RFC3339Nano}
	case 1:
		return FTime{time: t, format: format[0]}
	default:
		panic(errVarCount)
	}
}

// Equal determines if two FTimes are equal.
func (ft FTime) Equal(fTime FTime) bool {
	return ft.time.Equal(fTime.time) && ft.format == fTime.format
}

// Format returns the underlying format.
func (ft FTime) Format() string {
	return ft.format
}

// ParseFTime returns a new FTime parsed from the default time
// formats and any additional time formats.
func ParseFTime(timeStr string, additionalTimeFmts ...string) (FTime, error) {
	for i := 0; i < len(additionalTimeFmts); i++ {
		if time, err := time.Parse(additionalTimeFmts[i], timeStr); err == nil {
			return FTime{time: time, format: additionalTimeFmts[i]}, nil
		}
	}

	for i := 0; i < len(timeFmts); i++ {
		if time, err := time.Parse(timeFmts[i], timeStr); err == nil {
			return FTime{time: time, format: timeFmts[i]}, nil
		}
	}

	return FTime{}, errTimeFmt
}

// String returns a timestamp formatted by the underlying
// format.
func (ft FTime) String() string {
	return ft.time.Format(ft.format)
}

// Time returns the underlying time.Time.
func (ft FTime) Time() time.Time {
	return ft.time
}
