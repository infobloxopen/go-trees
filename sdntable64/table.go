package sdntable64

import (
	"fmt"
	"path/filepath"

	"github.com/infobloxopen/go-trees/udomain"
)

type Option func(o *options)

func WithPath(path string) Option {
	return func(o *options) {
		o.path = &path
	}
}

func WithLimit(limit int) Option {
	return func(o *options) {
		o.limit = limit
	}
}

func WithNormalizeLogger(f func(size, length int)) Option {
	return func(o *options) {
		o.log.normalize = f
	}
}

func WithFlushLogger(f func(size, from, to int)) Option {
	return func(o *options) {
		o.log.flush = f
	}
}

type loggers struct {
	normalize func(size, length int)
	flush     func(size, from, to int)
}

type options struct {
	path  *string
	limit int
	log   loggers
}

type Table64 struct {
	opts options
	root *uint64
	body [127]domains
}

func NewTable64(opts ...Option) *Table64 {
	t := new(Table64)
	for _, opt := range opts {
		opt(&t.opts)
	}

	if t.opts.path != nil {
		for i := range t.body {
			t.body[i] = makeDomainsWithPath(filepath.Join(*t.opts.path, fmt.Sprintf("%03d", i)))
		}

		if t.opts.limit < 10*1024*1024 {
			t.opts.limit = 100 * 1024 * 1024
		}
	}

	return t
}

func (t *Table64) InplaceInsert(k domain.Name, v uint64) {
	c := k.GetComparable()
	if len(c) < 1 {
		if t.root != nil {
			*t.root |= v
			return
		}

		t.root = &v
		return
	}

	i := len(c) - 1
	t.body[i] = t.body[i].inplaceInsert(c, v)
	t.flush()
}

func (t *Table64) RawInsert(k domain.Name, v uint64) {
	c := k.GetComparable()
	if len(c) < 1 {
		if t.root != nil {
			*t.root |= v
			return
		}

		t.root = &v
		return
	}

	i := len(c) - 1
	t.body[i] = t.body[i].rawInsert(c, v)
	t.flush()
}

func (t *Table64) Normalize() {
	for i, d := range t.body {
		t.body[i] = d.normalize(t.opts.log.normalize)
	}
}

func (t *Table64) Get(k domain.Name) (uint64, bool) {
	c := k.GetComparable()
	for len(c) > 0 {
		if v, ok := t.body[len(c)-1].get(c); ok {
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

	for i := range t.collectBigLists() {
		t.body[i] = t.body[i].flush(t.opts.log.flush, t.opts.log.normalize)
	}
}

func (t *Table64) collectBigLists() map[int]struct{} {
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
