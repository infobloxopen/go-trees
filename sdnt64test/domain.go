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
