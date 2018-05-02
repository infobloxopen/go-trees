package domain

import (
	"fmt"
	"math"
)

const (
	escRegular = iota
	escFirstChar
	escSecondDigit
	escThirdDigit
)

func markLabels(s string, offs []int) (int, error) {
	n := 0
	start := 0
	esc := escRegular
	var code int
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch esc {
		case escRegular:
			switch c {
			case '.':
				if start >= i {
					return 0, ErrEmptyLabel
				}

				if n >= len(offs) {
					return 0, ErrTooManyLabels
				}

				offs[n] = start
				n++

				start = i + 1

			case '\\':
				esc = escFirstChar
			}

		case escFirstChar:
			if c < '0' || c > '9' {
				esc = escRegular
			} else {
				code = int(c-'0') * 100
				if code > math.MaxUint8 {
					return 0, ErrInvalidEscape
				}

				esc = escSecondDigit
			}

		case escSecondDigit:
			if c < '0' || c > '9' {
				return 0, ErrInvalidEscape
			}

			code += int(c-'0') * 10
			if code > math.MaxUint8 {
				return 0, ErrInvalidEscape
			}

			esc = escThirdDigit

		case escThirdDigit:
			if c < '0' || c > '9' {
				return 0, ErrInvalidEscape
			}

			code += int(c - '0')
			if code > math.MaxUint8 {
				return 0, ErrInvalidEscape
			}

			esc = escRegular
		}
	}

	if esc != escRegular {
		return 0, ErrInvalidEscape
	}

	if start < len(s) {
		if n >= len(offs) {
			return 0, ErrTooManyLabels
		}

		offs[n] = start
		n++
	}

	return n, nil
}

func getLabel(s string, out []byte) (int, error) {
	j := 1
	esc := escRegular
	var code int

Loop:
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch esc {
		case escRegular:
			switch c {
			default:
				if c >= 'a' && c <= 'z' {
					c &= 0xdf
				}

				if j >= len(out) {
					return 0, ErrLabelTooLong
				}

				out[j] = c
				j++

			case '.':
				if i < len(s)-1 {
					panic(fmt.Errorf("unescaped dot at %d index before last character %d", i, len(s)-1))
				}

				break Loop

			case '\\':
				esc = escFirstChar
			}

		case escFirstChar:
			if c < '0' || c > '9' {
				if c >= 'a' && c <= 'z' {
					c &= 0xdf
				}

				if j >= len(out) {
					return 0, ErrLabelTooLong
				}

				out[j] = c
				j++

				esc = escRegular
			} else {
				code = int(c-'0') * 100
				if code > math.MaxUint8 {
					return 0, ErrInvalidEscape
				}

				esc = escSecondDigit
			}

		case escSecondDigit:
			if c < '0' || c > '9' {
				return 0, ErrInvalidEscape
			}

			code += int(c-'0') * 10
			if code > math.MaxUint8 {
				return 0, ErrInvalidEscape
			}

			esc = escThirdDigit

		case escThirdDigit:
			if c < '0' || c > '9' {
				return 0, ErrInvalidEscape
			}

			code += int(c - '0')
			if code > math.MaxUint8 {
				return 0, ErrInvalidEscape
			}

			c = byte(code)
			if c >= 'a' && c <= 'z' {
				c &= 0xdf
			}

			if j >= len(out) {
				return 0, ErrLabelTooLong
			}

			out[j] = c
			j++

			esc = escRegular
		}
	}

	if esc != escRegular {
		return 0, ErrInvalidEscape
	}

	out[0] = byte(j - 1)

	return j, nil
}
