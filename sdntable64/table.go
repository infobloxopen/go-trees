package sdntable64

import (
	"fmt"

	"github.com/infobloxopen/go-trees/udomain"
)

const (
	maxDomainLength = 127
	megaByte        = 1024 * 1024
)

type Table64 struct {
	opts  options
	root  uint64
	ready bool
	body  [maxDomainLength]domains
}

func NewTable64(opts ...Option) *Table64 {
	t := &Table64{
		ready: true,
	}

	for _, opt := range opts {
		opt(&t.opts)
	}

	if t.opts.path != nil {
		path := *t.opts.path
		for i := range t.body {
			t.body[i] = makeDomainsWithPath(path, i+1)
		}

		if t.opts.limit < 10*megaByte {
			t.opts.limit = 100 * megaByte
		}
	} else {
		for i := range t.body {
			t.body[i] = makeDomains()
		}
	}

	return t
}

func (t *Table64) InplaceInsert(k domain.Name, v uint64) {
	c := k.GetComparable()
	if len(c) > 0 {
		if !t.ready {
			panic(fmt.Errorf("table is not ready"))
		}

		i := len(c) - 1
		t.body[i] = t.body[i].inplaceInsert(c, v)
		t.flush()

		return
	}

	t.root = v
}

func (t *Table64) Append(k domain.Name, v uint64) {
	c := k.GetComparable()
	if len(c) > 0 {
		i := len(c) - 1
		t.body[i] = t.body[i].append(c, v)
		t.ready = false

		t.flush()
		return
	}

	t.root = v
}

func (t *Table64) Normalize() {
	for i, d := range t.body {
		if !d.ready {
			t.body[i] = d.normalize(t.opts.log.normalize)
		}
	}

	t.ready = true
}

func (t *Table64) Get(k domain.Name) uint64 {
	c := k.GetComparable()
	if len(c) > 0 && !t.ready {
		panic(fmt.Errorf("table is not ready"))
	}

	for len(c) > 0 {
		if v := t.body[len(c)-1].get(c, t.opts.log.read); v != 0 {
			return v
		}

		k = k.DropFirstLabel()
		c = k.GetComparable()
	}

	return t.root
}

func (t *Table64) Size() int {
	s := 0
	for _, d := range t.body {
		s += d.size()
	}

	return s
}

func (t *Table64) flush() {
	if t.opts.path == nil {
		return
	}

	for i := range t.makeFlushingList() {
		t.body[i] = t.body[i].flush(t.opts.log.flush, t.opts.log.normalize)
	}
}

func (t *Table64) makeFlushingList() map[int]struct{} {
	size := t.Size()
	if size < t.opts.limit {
		return nil
	}

	out := make(map[int]struct{}, len(t.body))

	prev := -1
	half := t.opts.limit / 2
	for prev != size && size > half {
		max := 0
		i := -1

		for j, d := range t.body {
			if _, ok := out[j]; !ok {
				s := d.size90()
				if s > max {
					max = s
					i = j
				}
			}
		}

		prev = size
		if i >= 0 {
			size -= max
			out[i] = struct{}{}
		}
	}

	return out
}
