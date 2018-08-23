package sdntable64

import "github.com/infobloxopen/go-trees/udomain"

type Table64 struct {
	root *uint64
	body [127]domains
}

func NewTable64() *Table64 {
	return new(Table64)
}

func (t *Table64) InplaceInsert(k domain.Name, v uint64) {
	c := k.GetComparable()
	if len(c) <= 0 {
		if t.root != nil {
			*t.root |= v
			return
		}

		t.root = &v
		return
	}

	t.body[len(c)] = t.body[len(c)].inplaceInsert(c, v)
}

func (t *Table64) Get(k domain.Name) (uint64, bool) {
	c := k.GetComparable()
	for len(c) > 0 {
		if v, ok := t.body[len(c)].get(c); ok {
			return v, ok
		}

		k = k.DropFirstLabel()
		c = k.GetComparable()
	}

	if t.root != nil {
		return *t.root, true
	}

	return 0, false
}
