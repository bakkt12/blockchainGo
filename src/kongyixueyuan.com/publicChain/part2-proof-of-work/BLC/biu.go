package BLC

import (
	"errors"
	"regexp"
)

const (
	zero  = byte('0')
	one   = byte('1')
	lsb   = byte('[')
	rsb   = byte(']')
	space = byte(' ')
)

var uint8arr [8]uint8
var ErrBadStringFormat = errors.New("bad string format")
var ErrEmptyString = errors.New("empty string")

func init() {
	uint8arr[0] = 128
	uint8arr[1] = 64
	uint8arr[2] = 32
	uint8arr[3] = 16
	uint8arr[4] = 8
	uint8arr[5] = 4
	uint8arr[6] = 2
	uint8arr[7] = 1
}
func appendBinaryString(bs []byte, b byte) []byte {
	var a byte
	for i := 0; i < 8; i++ {
		a = b
		b <<= 1
		b >>= 1
		switch a {
		case b:
			bs = append(bs, zero)
		default:
			bs = append(bs, one)
		}
		b <<= 1
	}
	return bs
}
func ByteToBinaryString(b byte) string {
	buf := make([]byte, 0, 8)
	buf = appendBinaryString(buf, b)
	return string(buf)
}
func BytesToBinaryString(bs []byte) string {
	l := len(bs)
	bl := l*8 + l + 1
	buf := make([]byte, 0, bl)
	buf = append(buf, lsb)
	for _, b := range bs {
		buf = appendBinaryString(buf, b)
		buf = append(buf, space)
	}
	buf[bl-1] = rsb
	return string(buf)
}

var rbDel = regexp.MustCompile(`[^01]`)

func BinaryStringToBytes(s string) (bs []byte) {
	if len(s) == 0 {
		panic(ErrEmptyString)
	}
	s = rbDel.ReplaceAllString(s, "")
	l := len(s)
	if l == 0 {
		panic(ErrBadStringFormat)
	}
	mo := l % 8
	l /= 8
	if mo != 0 {
		l++
	}
	bs = make([]byte, 0, l)
	mo = 8 - mo
	var n uint8
	for i, b := range []byte(s) {
		m := (i + mo) % 8
		switch b {
		case one:
			n += uint8arr[m]
		}
		if m == 7 {
			bs = append(bs, n)
			n = 0
		}
	}
	return
}
