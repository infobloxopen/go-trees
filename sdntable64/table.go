package sdntable64

import "github.com/infobloxopen/go-trees/udomain"

type Table64 struct {
	root *uint64
	body [2080][]uint64
}

func NewTable64() *Table64 {
	return new(Table64)
}

func (t *Table64) InplaceInsert(k domain.Name, v uint64) {
	n := k.GetLabelCount()
	m := domain.MaxLabels - n
	if m >= domain.MaxLabels {
		t.root = &v
		return
	}

	c := k.GetComparable()
	s := len(c) + 1

	i := getDomainNameIndex(m) + len(c) - n
	a := t.body[i]
	if a == nil {
		a = make([]uint64, s)
		copy(a, c)
		a[len(a)-1] = v
	} else {
		end := len(a) / s

		left := 0
		right := end
		for {
			mid := (right + left) / 2
			if k.Less(domain.MakeNameFromSlice(a[s*mid : s*(mid+1)-1])) {
				right = mid
				if left == right {
					break
				}
			} else if left == mid {
				left = right
				break
			} else {
				left = mid
			}
		}

		a = append(append(a, c...), v)
		if left < end {
			copy(a[s*(left+1):], a[s*left:])
			copy(a[s*left:], c)
			a[s*left+len(c)] = v
		}
	}

	t.body[i] = a
}

func (t *Table64) Get(k domain.Name) (uint64, bool) {
	n := k.GetLabelCount()
	m := domain.MaxLabels - n
	for m < domain.MaxLabels {
		c := k.GetComparable()
		s := len(c) + 1

		if a := t.body[getDomainNameIndex(m)+len(c)-n]; a != nil {
			end := len(a) / s

			left := 0
			right := end
			for {
				mid := (right + left) / 2
				mdn := domain.MakeNameFromSlice(a[s*mid : s*(mid+1)-1])
				if k.Less(mdn) {
					right = mid
					if left == right {
						break
					}
				} else if left == mid {
					if !mdn.Less(k) {
						return a[s*(mid+1)-1], true
					}

					break
				} else {
					left = mid
				}
			}
		}

		k = k.DropFirstLabel()
		n = k.GetLabelCount()
		m = domain.MaxLabels - n
	}

	if t.root != nil {
		return *t.root, true
	}

	return 0, false
}

func getDomainNameIndex(n int) int {
	x := n / 4
	y := n % 4
	return (2*x + y) * (x + 1)
}
