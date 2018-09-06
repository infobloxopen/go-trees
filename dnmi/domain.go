package main

import "errors"

var errDomainOverflow = errors.New("overflow")

func prevDomain(a []int64) ([]int64, error) {
	out := make([]int64, len(a))
	copy(out, a)

	if domainDec(out) {
		return out, errDomainOverflow
	}

	return out, nil
}

func domainDec(out []int64) bool {
	if len(out) > 0 {
		if cf := domainDec(out[int(out[0]&0xf):]); cf {
			return domainLabelDec(out)
		}

		return false
	}

	return true
}

func domainLabelDec(out []int64) bool {
	n := int(out[0] & 0xf)
	m := int(out[0] & 0xf0 >> 4)
	if m <= 0 {
		m = 8
	}

	for i := n - 1; i >= 0; i-- {
		j := 0
		if i == 0 {
			j++
		}

		for j < m {
			shift := uint(j * 8)
			c := uint64(out[i]>>shift) & 0xff
			cf := false
			if c > 'Z' {
				c = 'Z'
			} else if c > 'A' {
				c -= 1
			} else if c > '9' {
				c = '9'
			} else if c > '0' {
				c -= 1
			} else {
				c = 'Z'
				cf = true
			}

			out[i] = int64(uint64(out[i])&masks[j] | c<<shift)

			if !cf {
				return false
			}

			j++
		}

		m = 8
	}

	return true
}

func nextDomain(a []int64) ([]int64, error) {
	out := make([]int64, len(a))
	copy(out, a)

	if domainInc(out) {
		return out, errDomainOverflow
	}

	return out, nil
}

func domainInc(out []int64) bool {
	if len(out) > 0 {
		if cf := domainInc(out[int(out[0]&0xf):]); cf {
			return domainLabelInc(out)
		}

		return false
	}

	return true
}

func domainLabelInc(out []int64) bool {
	n := int(out[0] & 0xf)
	m := int(out[0] & 0xf0 >> 4)
	if m <= 0 {
		m = 8
	}

	for i := n - 1; i >= 0; i-- {
		j := 0
		if i == 0 {
			j++
		}

		for j < m {
			shift := uint(j * 8)
			c := uint64(out[i]>>shift) & 0xff
			cf := false
			if c < '0' {
				c = '0'
			} else if c < '9' {
				c++
			} else if c < 'A' {
				c = 'A'
			} else if c < 'Z' {
				c++
			} else {
				c = '0'
				cf = true
			}

			out[i] = int64(uint64(out[i])&masks[j] | c<<shift)

			if !cf {
				return false
			}

			j++
		}

		m = 8
	}

	return true
}

var masks = []uint64{
	0xffffffffffffff00,
	0xffffffffffff00ff,
	0xffffffffff00ffff,
	0xffffffff00ffffff,
	0xffffff00ffffffff,
	0xffff00ffffffffff,
	0xff00ffffffffffff,
	0x00ffffffffffffff,
}
