package main

type sortable struct {
	i []int
	a []int64
}

func makeSortable(in []int64, n int) sortable {
	out := sortable{
		i: make([]int, len(in)/n),
		a: in,
	}

	for i := range out.i {
		out.i[i] = i
	}

	return out
}

func (s sortable) getSorted() []int64 {
	if len(s.i) < 2 {
		return s.a
	}

	m := len(s.a) / len(s.i)

	n := 1
	for i, j := range s.i[1:] {
		a := s.i[i] * m
		if compare(s.a[a:a+m], s.a[j*m:]) != 0 {
			n++
		}
	}

	out := make([]int64, n*m)
	copy(out[:m], s.a[s.i[0]*m:])

	a := 0
	for _, j := range s.i[1:] {
		b := j * m
		if compare(out[a:a+m], s.a[b:]) != 0 {
			a += m
			copy(out[a:a+m], s.a[b:])
		}
	}

	return out
}

func (s sortable) Len() int      { return len(s.i) }
func (s sortable) Swap(i, j int) { s.i[i], s.i[j] = s.i[j], s.i[i] }
func (s sortable) Less(i, j int) bool {
	m := len(s.a) / len(s.i)
	i, j = s.i[i]*m, s.i[j]*m
	return compare(s.a[i:i+m], s.a[j:]) < 0
}

func compare(a, b []int64) int64 {
	for i, n := range a {
		if d := n - b[i]; d != 0 {
			return d
		}
	}

	return 0
}
