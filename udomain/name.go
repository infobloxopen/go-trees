// Package domain provides functions to parse and handle domain names.
package domain

import (
	"encoding/binary"
	"errors"
	"math"
)

var (
	// ErrLabelTooLong is returned when one of domain labels has more than 63 characters.
	ErrLabelTooLong = errors.New("label too long")
	// ErrEmptyLabel indicates that domain name contains empty label.
	ErrEmptyLabel = errors.New("empty label")
	// ErrNameTooLong is returned when domain name has more than 255 characters.
	ErrNameTooLong = errors.New("name too long")
	// ErrInvalidEscape is returned for invalid escape sequence.
	ErrInvalidEscape = errors.New("invalid escape sequence")
	// ErrInvalidLabelSize
	ErrInvalidLabelSize = errors.New("invalid label size")
	// ErrInvalidCharacter
	ErrInvalidCharacter = errors.New("invalid character")
)

// Name is a structure which represents domain name.
type Name struct {
	h string
	n int
	c []int64
}

const (
	// MaxName is maximum number of bytes for whole domain name.
	MaxName = 255
	// MaxLabels is maximum number of labels domain name can consist of.
	MaxLabels = MaxName / 2
	// MaxLabel is maximim number of bytes for single label.
	MaxLabel = 63

	uInt64Size = 8
)

const (
	escRegular = iota
	escFirstChar
	escSecondDigit
	escThirdDigit
)

func MakeNameFromString(s string) (Name, error) {
	return MakeNameFromStringWithBuffer(s, make([]int64, 0, MaxLabels))
}

func MakeNameFromStringWithBuffer(s string, buf []int64) (Name, error) {
	out := Name{h: s}
	if len(s) <= 0 || s == "." {
		return out, nil
	}

	var (
		fragment [uInt64Size]byte
		zeros    [uInt64Size]byte
		label    [(MaxLabel + 1) / uInt64Size]int64
	)

	n := 0
	j := 1
	count := 1
	esc := escRegular
	code := 0
	for i := 0; i < len(s); i++ {
		b := s[i]
		switch esc {
		case escRegular:
			switch b {
			default:
				if b >= 'a' && b <= 'z' {
					b &= 0xdf
				}

				fragment[j] = b
				j++

				count++
				if count > MaxName {
					return out, ErrNameTooLong
				}

				if j >= len(fragment) {
					if n >= len(label) {
						return out, ErrLabelTooLong
					}

					label[n] = int64(binary.LittleEndian.Uint64(fragment[:]))
					n++
					j = 0
				}

			case '.':
				if n < 1 && j <= 1 {
					return out, ErrEmptyLabel
				}

				if j > 0 {
					if n >= len(label) {
						return out, ErrLabelTooLong
					}

					copy(fragment[j:], zeros[:])

					label[n] = int64(binary.LittleEndian.Uint64(fragment[:]))
					n++
				}

				label[0] = int64(uint64(label[0]) | uint64(n) | uint64(j<<4))

				j = 1
				count++

				buf = append(buf, label[:n]...)
				out.n++

				n = 0
				fragment[0] = 0

			case '\\':
				esc = escFirstChar
			}

		case escFirstChar:
			if b < '0' || b > '9' {
				esc = escRegular

				if b >= 'a' && b <= 'z' {
					b &= 0xdf
				}

				fragment[j] = b
				j++

				count++
				if count > MaxName {
					return out, ErrNameTooLong
				}

				if j >= len(fragment) {
					if n >= len(label) {
						return out, ErrLabelTooLong
					}

					label[n] = int64(binary.LittleEndian.Uint64(fragment[:]))
					n++
					j = 0
				}

				break
			}

			esc = escSecondDigit

			code = int(b-'0') * 100
			if code > math.MaxUint8 {
				return out, ErrInvalidEscape
			}

		case escSecondDigit:
			if b < '0' || b > '9' {
				return out, ErrInvalidEscape
			}

			esc = escThirdDigit

			code += int(b-'0') * 10
			if code > math.MaxUint8 {
				return out, ErrInvalidEscape
			}

		case escThirdDigit:
			if b < '0' || b > '9' {
				return out, ErrInvalidEscape
			}

			esc = escRegular

			code += int(b - '0')
			if code > math.MaxUint8 {
				return out, ErrInvalidEscape
			}

			b = byte(code)
			if b >= 'a' && b <= 'z' {
				b &= 0xdf
			}

			fragment[j] = b
			j++

			count++
			if count > MaxName {
				return out, ErrNameTooLong
			}

			if j >= len(fragment) {
				if n >= len(label) {
					return out, ErrLabelTooLong
				}

				label[n] = int64(binary.LittleEndian.Uint64(fragment[:]))
				n++
				j = 0
			}
		}
	}

	if esc != escRegular {
		return out, ErrInvalidEscape
	}

	if n > 0 || j > 1 {
		if n >= len(label) {
			return out, ErrLabelTooLong
		}

		if j > 0 {
			copy(fragment[j:], zeros[:])

			label[n] = int64(binary.LittleEndian.Uint64(fragment[:]))
			n++
		}

		label[0] = int64(uint64(label[0]) | uint64(n) | uint64(j<<4))

		buf = append(buf, label[:n]...)
		out.n++
	}

	out.c = buf
	return out, nil
}

func MakeNameFromSlice(s []int64) (Name, error) {
	out := Name{c: s}

	b := []byte{}
	for len(s) > 0 {
		n := int(s[0] & 0xf)
		if n < 1 || n > 8 || n > len(s) {
			return out, ErrInvalidLabelSize
		}

		j := 1
		for i := 0; i < n; i++ {
			m := 8
			if i == n-1 {
				m = int(s[0] & 0xf0 >> 4)
				if m <= 0 {
					m = 8
				}

				if m > 8 {
					return out, ErrInvalidLabelSize
				}
			}

			for j < m {
				c := byte((s[i] >> uint(8*j)) & 0xff)
				if c >= 'a' && c <= 'z' {
					return out, ErrInvalidCharacter
				}

				if c == '-' || c == '_' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' {
					if c >= 'A' && c <= 'Z' {
						c |= 0x20
					}
					b = append(b, c)
				} else if c >= '!' && c <= '~' {
					b = append(b, '\\', c)
				} else {
					b = append(b, '\\', c/100+'0', c%100/10+'0', c%10+'0')
				}
				j++
			}

			j = 0
		}

		s = s[n:]
		if len(s) > 0 {
			b = append(b, '.')
		}
	}

	out.h = string(b)

	return out, nil
}

func (n Name) String() string {
	return n.h
}

func (n Name) DropFirstLabel() Name {
	if n.n > 1 && len(n.c) > 0 {
		if i := n.c[0] & 7; int(i) <= len(n.c) {
			return Name{
				n: n.n - 1,
				c: n.c[i:],
			}
		}
	}

	return Name{}
}

func (n Name) GetComparable() []int64 {
	return n.c
}

func (n Name) Less(other Name) bool {
	if len(n.c) < len(other.c) {
		return true
	}

	for i, u := range n.c {
		if r := u - other.c[i]; r != 0 {
			return r < 0
		}
	}

	return false
}
