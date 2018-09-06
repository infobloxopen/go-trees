package main

import (
	"math"

	"github.com/infobloxopen/go-trees/udomain"
)

const (
	escRegular = iota
	escFirstChar
	escSecondDigit
	escThirdDigit
)

func domainUppercase(s string) (string, error) {
	b := make([]byte, len(s))
	n := 0
	esc := escRegular
	code := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch esc {
		case escRegular:
			switch c {
			default:
				if c >= 'a' && c <= 'z' {
					c &= 0xdf
				}

				b[n] = c
				n++

			case '\\':
				esc = escFirstChar
			}

		case escFirstChar:
			if c < '0' || c > '9' {
				if c >= 'a' && c <= 'z' {
					c &= 0xdf
				}

				if c == '.' || c == '\\' {
					b[n] = '\\'
					b[n+1] = c
					n += 2
				} else {
					b[n] = c
					n++
				}

				esc = escRegular
				break
			}

			code = int(c-'0') * 100
			if code > math.MaxUint8 {
				return s, domain.ErrInvalidEscape
			}

			esc = escSecondDigit

		case escSecondDigit:
			if c < '0' || c > '9' {
				return s, domain.ErrInvalidEscape
			}

			code += int(c-'0') * 10
			if code > math.MaxUint8 {
				return s, domain.ErrInvalidEscape
			}

			esc = escThirdDigit

		case escThirdDigit:
			if c < '0' || c > '9' {
				return s, domain.ErrInvalidEscape
			}

			code += int(c - '0')
			if code > math.MaxUint8 {
				return s, domain.ErrInvalidEscape
			}

			c = byte(code)
			if c >= 'a' && c <= 'z' {
				c &= 0xdf
			}

			if c == '.' || c == '\\' {
				b[n] = '\\'
				b[n+1] = c
				n += 2
			} else {
				b[n] = c
				n++
			}

			esc = escRegular
		}
	}

	return string(b[:n]), nil
}

func dropLabel(s string) (string, error) {
	esc := escRegular
	code := 0
	for i := 0; i < len(s); i++ {
		b := s[i]
		switch esc {
		case escRegular:
			switch b {
			case '.':
				return s[i+1:], nil

			case '\\':
				esc = escFirstChar
			}

		case escFirstChar:
			if b < '0' || b > '9' {
				esc = escRegular
				break
			}

			code = int(b-'0') * 100
			if code > math.MaxUint8 {
				return s, domain.ErrInvalidEscape
			}

			esc = escSecondDigit

		case escSecondDigit:
			if b < '0' || b > '9' {
				return s, domain.ErrInvalidEscape
			}

			code += int(b-'0') * 10
			if code > math.MaxUint8 {
				return s, domain.ErrInvalidEscape
			}

			esc = escThirdDigit

		case escThirdDigit:
			if b < '0' || b > '9' {
				return s, domain.ErrInvalidEscape
			}

			code += int(b - '0')
			if code > math.MaxUint8 {
				return s, domain.ErrInvalidEscape
			}

			esc = escRegular
		}
	}

	return "", nil
}
