package sdntable64

import "io"

type wRun struct {
	buf  run
	w    io.Writer
	blks int
	rem  int
	k    []int64
	i    int
	n    int
}

func (d domains) makeWriterRun(w io.Writer) *wRun {
	n := len(d.data.k) / len(d.data.v)
	buf := makeRunForBlock(n)
	return &wRun{
		buf: buf,
		w:   w,
		k:   buf.k,
		n:   n,
	}
}

func (r *wRun) put(k []int64, v uint64) error {
	copy(r.k, k)
	r.buf.v[r.i] = v

	r.k = r.k[r.n:]
	r.i++

	if len(r.k) <= 0 {
		if err := r.buf.write(r.w); err != nil {
			return err
		}

		r.k = r.buf.k
		r.i = 0

		r.blks++
	}

	return nil
}

func (r *wRun) flush() error {
	if r.i > 0 {
		if err := r.buf.truncate(r.i).write(r.w); err != nil {
			return err
		}

		r.rem = r.i
	}

	return nil
}
