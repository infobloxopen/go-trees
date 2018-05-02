package domain

import "fmt"

const (
	escRegular = iota
	escFirstChar
)

func markLabels(s string, offs []int) (int, error) {
	n := 0
	start := 0
	esc := escRegular
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
			esc = escRegular

			if c >= 'a' && c <= 'z' {
				c &= 0xdf
			}

			if j >= len(out) {
				return 0, ErrLabelTooLong
			}

			out[j] = c
			j++
		}
	}

	if esc != escRegular {
		return 0, ErrInvalidEscape
	}

	out[0] = byte(j - 1)

	return j, nil
}
