package fmtstate

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

type CustomState struct {
	width     int
	precision int
	flags     string
	buf       []byte
}

func (cs *CustomState) Write(b []byte) (n int, err error) {
	cs.buf = append(cs.buf, b...)
	return len(b), nil
}

func (cs *CustomState) Width() (wid int, ok bool) {
	return cs.width, cs.width > 0
}

func (cs *CustomState) Precision() (prec int, ok bool) {
	return cs.precision, cs.precision > 0
}

func (cs *CustomState) Flag(c int) bool {
	return strings.ContainsRune(cs.flags, rune(c))
}

func (cs *CustomState) String() string {
	return string(cs.buf)
}

func NewCustomState(width int, precision int, flags string) *CustomState {
	return &CustomState{
		width:     width,
		precision: precision,
		flags:     flags,
		buf:       make([]byte, 0, 64), // initial capacity of 64 bytes
	}
}

func (cs *CustomState) WriteString(s string) (n int, err error) {
	cs.buf = append(cs.buf, s...)
	return len(s), nil
}

// WriteRune is a helper method to write single runes
func (cs *CustomState) WriteRune(r rune) (n int, err error) {
	if r < utf8.RuneSelf {
		cs.buf = append(cs.buf, byte(r))
		return 1, nil
	}
	var buf [utf8.UTFMax]byte
	n = utf8.EncodeRune(buf[:], r)
	cs.buf = append(cs.buf, buf[:n]...)
	return n, nil
}

// WriteInt is a helper method to write integers
func (cs *CustomState) WriteInt(i int64) (n int, err error) {
	return cs.WriteString(strconv.FormatInt(i, 10))
}
