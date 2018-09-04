package sdntable64

import "io"

type rRun struct {
	buf  run
	r    io.Reader
	blks int
	rem  int
	k    []int64
	i    int
	n    int
}

func (d domains) makeReaderRun(r io.Reader) *rRun {
	n := len(d.data.k) / len(d.data.v)
	buf := makeRunForBlock(n)
	return &rRun{
		buf:  buf,
		r:    r,
		blks: d.blks,
		rem:  d.rem,
		n:    n,
	}
}

func (r *rRun) isEmpty() bool {
	return r.blks <= 0 && r.rem <= 0 && len(r.k) <= 0
}

func (r *rRun) next() ([]int64, uint64, error) {
	if len(r.k) <= 0 {
		if r.r == nil {
			return nil, 0, nil
		}

		if r.blks > 0 {
			r.blks--
		} else if r.rem > 0 {
			r.buf = r.buf.truncate(r.rem)
			r.rem = 0
		} else {
			return nil, 0, nil
		}

		if err := r.buf.read(r.r); err != nil {
			return nil, 0, err
		}

		r.k = r.buf.k
		r.i = 0
	}

	k := r.k[:r.n]
	v := r.buf.v[r.i]

	r.k = r.k[r.n:]
	r.i++

	return k, v, nil
}
