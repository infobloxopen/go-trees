package sdntable64

import "sort"

type indexedRun struct {
	i []int
	k []int64
	v []uint64
	d []uint32
}

func (r run) makeIndexedRun() indexedRun {
	idx := make([]int, len(r.v))
	for i := range idx {
		idx[i] = i
	}

	return indexedRun{
		i: idx,
		k: r.k,
		v: r.v,
		d: r.i,
	}
}

func (r indexedRun) makeRun() run {
	m := len(r.v)
	n := len(r.k) / m
	for k, i := range r.i[:len(r.i)-1] {
		j := r.i[k+1]
		a, b := n*i, n*j
		if compare(r.k[a:a+n], r.k[b:b+n]) == 0 {
			m--
		}
	}

	k := make([]int64, m*n)
	v := make([]uint64, m)

	var d []uint32
	if len(r.d) > 0 {
		d = make([]uint32, m)
	}

	copy(k[:n], r.k[n*r.i[0]:])
	v[0] = r.v[r.i[0]]
	if d != nil {
		d[0] = r.d[r.i[0]]
	}

	p := 0
	for i, j := range r.i[1:] {
		a, b := n*r.i[i], n*j
		if compare(r.k[a:a+n], r.k[b:b+n]) != 0 {
			p++
			copy(k[n*p:], r.k[b:b+n])
		}

		v[p] = r.v[j]
		if d != nil {
			d[p] = r.d[j]
		}
	}

	return run{
		k: k,
		v: v,
		i: d,
	}
}

func (r indexedRun) sort() {
	sort.Sort(r)
}

func (r indexedRun) Len() int      { return len(r.i) }
func (r indexedRun) Swap(i, j int) { r.i[i], r.i[j] = r.i[j], r.i[i] }
func (r indexedRun) Less(i, j int) bool {
	n := len(r.k) / len(r.v)
	i, j = n*r.i[i], n*r.i[j]
	if d := compare(r.k[i:i+n], r.k[j:j+n]); d != 0 {
		return d < 0
	}

	return i < j
}

func compare(a, b []int64) int64 {
	for i, n := range a {
		if d := n - b[i]; d != 0 {
			return d
		}
	}

	return 0
}
