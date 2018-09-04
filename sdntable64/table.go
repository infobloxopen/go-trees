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

type chDomains struct {
	m int
	d domains
	g *Getter
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
			t.body[i] = makeDomainsWithPath(path)
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

func (t *Table64) Close() error {
	for _, d := range t.body {
		d.close()
	}

	return nil
}

func (t *Table64) Insert(k domain.Name, v uint64) (*Table64, []*Getter) {
	c := k.GetComparable()
	if len(c) > 0 {
		if !t.ready {
			panic(fmt.Errorf("table is not ready"))
		}

		i := len(c) - 1
		data, ok := t.body[i].insert(c, v)
		if !ok {
			return t, nil
		}

		out := &Table64{
			opts:  t.opts,
			root:  v,
			ready: t.ready,
			body:  t.body,
		}

		out.body[i] = data
		flushed := out.flush()
		if len(flushed) <= 0 {
			return out, nil
		}

		g := make([]*Getter, 0, len(flushed))
		for _, d := range flushed {
			out.body[d.m-1] = d.d
			if d.g != nil {
				g = append(g, d.g)
			}
		}

		return out, g
	}

	return &Table64{
		opts:  t.opts,
		root:  v,
		ready: t.ready,
		body:  t.body,
	}, nil
}

func (t *Table64) Delete(k domain.Name) (*Table64, []*Getter) {
	return t.Insert(k, 0)
}

func (t *Table64) InplaceInsert(k domain.Name, v uint64) {
	c := k.GetComparable()
	if len(c) > 0 {
		if !t.ready {
			panic(fmt.Errorf("table is not ready"))
		}

		i := len(c) - 1
		t.body[i] = t.body[i].inplaceInsert(c, v)
		for _, d := range t.flush() {
			t.body[d.m-1] = d.d
			if d.g != nil {
				if err := d.g.Stop(); err != nil {
					panic(fmt.Errorf("can't stop getter %q for subarray %d", d.g.path, d.m))
				}
			}
		}

		return
	}

	t.root = v
}

func (t *Table64) Append(k domain.Name, v uint64) (*Table64, []*Getter) {
	out := &Table64{
		opts:  t.opts,
		root:  t.root,
		ready: t.ready,
		body:  t.body,
	}

	c := k.GetComparable()
	if len(c) > 0 {
		i := len(c) - 1
		out.body[i] = t.body[i].append(c, v)
		out.ready = false

		flushed := out.flush()
		if len(flushed) <= 0 {
			return out, nil
		}

		g := make([]*Getter, 0, len(flushed))
		for _, d := range flushed {
			out.body[d.m-1] = d.d
			if d.g != nil {
				g = append(g, d.g)
			}
		}

		return out, g
	}

	out.root = v
	return out, nil
}

func (t *Table64) Normalize() *Table64 {
	out := &Table64{
		opts:  t.opts,
		root:  t.root,
		ready: t.ready,
	}

	for i, d := range t.body {
		if d.ready {
			out.body[i] = d
		} else {
			out.body[i] = d.normalize(t.opts.log.normalize)
		}
	}

	out.ready = true
	return out
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

func (t *Table64) flush() []chDomains {
	if t.opts.path == nil {
		return nil
	}

	idx := t.makeFlushingList()
	if len(idx) <= 0 {
		return nil
	}

	ch := make([]chDomains, 0, len(idx))
	for i := range idx {
		d := t.body[i]
		d, g := d.flush(t.opts.log.flush, t.opts.log.read, t.opts.log.normalize)

		ch = append(ch, chDomains{
			m: i + 1,
			d: d,
			g: g,
		})
	}

	return ch
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
