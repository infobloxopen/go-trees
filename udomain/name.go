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
)

// Name is a structure which represents domain name.
type Name struct {
	h string
	n int
	c []uint64
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
	out := Name{h: s}
	if len(s) <= 0 || s == "." {
		return out, nil
	}

	var (
		fragment [uInt64Size]byte
		zeros    [uInt64Size]byte
		label    [(MaxLabel + 1) / uInt64Size]uint64
		name     [MaxLabels]uint64
	)

	n := 0
	j := 1
	k := 0
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

					label[n] = binary.LittleEndian.Uint64(fragment[:])
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

					label[n] = binary.LittleEndian.Uint64(fragment[:])
					n++
				}

				label[0] |= uint64(n) | uint64(j<<4)

				j = 1
				count++

				copy(name[k:], label[:n])
				out.n++
				k += n

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

					label[n] = binary.LittleEndian.Uint64(fragment[:])
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

				label[n] = binary.LittleEndian.Uint64(fragment[:])
				n++
				j = 0
			}
		}
	}

	if esc != escRegular {
		return out, ErrInvalidEscape
	}

	if n > 0 && j > 0 || j > 1 {
		if n >= len(label) {
			return out, ErrLabelTooLong
		}

		copy(fragment[j:], zeros[:])

		label[n] = binary.LittleEndian.Uint64(fragment[:])
		n++

		label[0] |= uint64(n) | uint64(j<<4)

		copy(name[k:], label[:n])
		out.n++
		k += n
	}

	out.c = name[:k]
	return out, nil
}

func MakeNameFromSlice(s []uint64) Name {
	return Name{c: s}
}

func (n Name) String() string {
	return n.h
}

func (n Name) DropFirstLabel() Name {
	if n.n > 1 && len(n.c) > 0 {
		return Name{
			n: n.n - 1,
			c: n.c[n.c[0]&7:],
		}
	}

	return Name{}
}

func (n Name) GetLabelCount() int {
	return n.n
}

func (n Name) GetComparable() []uint64 {
	return n.c
}

func (n Name) Less(other Name) bool {
	if len(n.c) < len(other.c) {
		return true
	}

	if len(n.c) == len(other.c) {
		for i, a := range n.c {
			b := other.c[i]
			if a < b {
				return true
			}

			if a > b {
				break
			}
		}
	}

	return false
}
