package sdntable64

import (
	"fmt"
	"math"
	"os"
	"sort"
)

type domains struct {
	path   *string
	idx    []uint32
	keys   []int64
	values []uint64
	blocks int
	rem    int
	dIdx   []uint32
}

func makeDomainsWithPath(path string) domains {
	return domains{
		path: &path,
	}
}

func (d domains) inplaceInsert(k []int64, v uint64) domains {
	if len(d.idx) <= 0 {
		keys := make([]int64, len(k))
		copy(keys, k)

		return domains{
			path:   d.path,
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
	if d.dIdx != nil {
		d.dIdx = append(d.dIdx, math.MaxUint32)
	}

	return d
}

func (d domains) rawInsert(k []int64, v uint64) domains {
	if len(d.idx) <= 0 {
		keys := make([]int64, len(k))
		copy(keys, k)

		return domains{
			path:   d.path,
			idx:    []uint32{0},
			keys:   keys,
			values: []uint64{v},
		}
	}

	end := len(d.idx)
	if end > math.MaxUint32 {
		panic(fmt.Errorf("table overflow"))
	}

	d.idx = append(d.idx, uint32(end))
	d.keys = append(d.keys, k...)
	d.values = append(d.values, v)
	if d.dIdx != nil {
		d.dIdx = append(d.dIdx, math.MaxUint32)
	}

	return d
}

func (d domains) normalize(log func(size, length int)) domains {
	if len(d.idx) <= 0 {
		return d
	}

	n := len(d.keys) / len(d.idx)
	if log != nil {
		log(n, len(d.idx))
	}

	if len(d.idx) < 2 {
		return d
	}

	sort.Sort(d)

	c := 0
	for k, a := range d.idx[:len(d.idx)-1] {
		b := d.idx[k+1]
		i, j := int(a), int(b)
		if compare(d.keys[n*i:n*(i+1)], d.keys[n*j:n*(j+1)]) == 0 {
			c++
		}
	}

	if c > 0 {
		idx := make([]uint32, len(d.idx)-c)
		keys := make([]int64, len(idx)*n)
		values := make([]uint64, len(idx))

		j := int(d.idx[0])
		idx[0] = 0
		copy(keys[:n], d.keys[n*j:])
		values[0] = d.values[j]

		m := 1

		for k, a := range d.idx[:len(d.idx)-1] {
			b := d.idx[k+1]
			i, j := int(a), int(b)
			if compare(d.keys[n*i:n*(i+1)], d.keys[n*j:n*(j+1)]) != 0 {
				idx[m] = uint32(m)
				copy(keys[n*m:n*(m+1)], d.keys[n*j:])
				values[m] = d.values[j]
				m++
			}
		}

		d.idx = idx
		d.keys = keys
		d.values = values
	}

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

func (d domains) size() int {
	return len(d.idx)*4 + (len(d.keys)+len(d.values))*8
}

func (d domains) size90() int {
	if d.path == nil || len(d.idx) <= 0 {
		return 0
	}

	return len(d.idx) * 90 / 100 * (12 + 8*len(d.keys)/len(d.idx))
}

func (d domains) flush(logFlush func(size, fromLen, toLen int), logNorm func(size, length int)) domains {
	if d.path == nil || len(d.idx) <= 0 {
		return d
	}

	n := len(d.idx) / 10
	m := len(d.keys) / len(d.idx)

	if logFlush != nil {
		logFlush(m, len(d.idx), n)
	}

	d = d.normalize(logNorm)

	var path string
	if len(d.dIdx) > 0 {
		dIdx, blks, rem, tmp, err := d.merge()
		if err != nil {
			os.Remove(path)
			panic(err)
		}

		d.dIdx = dIdx
		d.blocks = blks
		d.rem = rem
		path = tmp
	} else {
		blks, rem, tmp, err := d.writeAll()
		if err != nil {
			os.Remove(path)
			panic(err)
		}

		d.blocks = blks
		d.rem = rem
		path = tmp
	}

	if err := os.Rename(path, *d.path); err != nil {
		os.Remove(path)
		panic(err)
	}

	idx := make([]uint32, n)
	keys := make([]int64, m*n)
	dstKeys := keys
	srcKeys := d.keys
	values := make([]uint64, n)
	dIdx := make([]uint32, n)

	j := 0
	for i := 0; i < n; i++ {
		idx[i] = uint32(i)
		copy(dstKeys[:m], srcKeys)
		values[i] = d.values[j]
		if len(d.dIdx) > 0 {
			dIdx[i] = d.dIdx[j]
		} else {
			dIdx[i] = uint32(j)
		}

		j += 10
		dstKeys = dstKeys[m:]
		srcKeys = srcKeys[10*m:]
	}

	return domains{
		path:   d.path,
		idx:    idx,
		keys:   keys,
		values: values,
		blocks: d.blocks,
		rem:    d.rem,
		dIdx:   d.dIdx,
	}
}

func (d domains) Len() int      { return len(d.idx) }
func (d domains) Swap(i, j int) { d.idx[i], d.idx[j] = d.idx[j], d.idx[i] }
func (d domains) Less(i, j int) bool {
	n := len(d.keys) / len(d.idx)
	return compare(d.keys[n*i:n*(i+1)], d.keys[n*j:n*(j+1)]) < 0
}

func compare(a, b []int64) int64 {
	for i, u := range a {
		if r := u - b[i]; r != 0 {
			return r
		}
	}

	return 0
}
