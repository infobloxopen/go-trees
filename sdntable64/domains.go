package sdntable64

import (
	"fmt"
	"math"
)

type domains struct {
	idx    []uint32
	keys   []int64
	values []uint64
}

func (d domains) inplaceInsert(k []int64, v uint64) domains {
	if len(d.idx) <= 0 {
		keys := make([]int64, len(k))
		copy(keys, k)

		return domains{
			idx:    []uint32{0},
			keys:   keys,
			values: []uint64{v},
		}
	}

	end := len(d.idx)
	left := 0
	right := end
	for {
		m := (left + right) / 2
		i := int(d.idx[m])
		r := compare(k, d.keys[i*len(k):(i+1)*len(k)])
		if r == 0 {
			d.values[i] |= v
			return d
		}

		if r < 0 {
			right = m
			if left == right {
				break
			}
		} else {
			if left == m {
				left = right
				break
			}

			left = m
		}
	}

	if end > math.MaxUint32 {
		panic(fmt.Errorf("table overflow"))
	}

	i := uint32(end)
	d.idx = append(d.idx, i)
	if left < end {
		copy(d.idx[left+1:], d.idx[left:end])
		d.idx[left] = i
	}

	d.keys = append(d.keys, k...)
	d.values = append(d.values, v)

	return d
}

func (d domains) get(k []int64) (uint64, bool) {
	if len(d.idx) > 0 {
		left := 0
		right := len(d.idx)
		for {
			m := (left + right) / 2
			i := int(d.idx[m])
			r := compare(k, d.keys[i*len(k):(i+1)*len(k)])
			if r == 0 {
				return d.values[i], true
			}

			if left == m {
				break
			}

			if r < 0 {
				right = m
			} else {
				left = m
			}
		}
	}

	return 0, false
}

func compare(a, b []int64) int64 {
	for i, u := range a {
		if r := u - b[i]; r != 0 {
			return r
		}
	}

	return 0
}
