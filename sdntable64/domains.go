package sdntable64

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type domains struct {
	dir    *string
	getter *Getter
	ready  bool
	data   run
	blks   int
	rem    int
}

func makeDomains() domains {
	return domains{
		ready: true,
	}
}

func makeDomainsWithPath(path string) domains {
	d := makeDomains()
	d.dir = &path

	return d
}

func (d domains) close() error {
	if d.getter != nil {
		return d.getter.Stop()
	}

	return nil
}

func (d domains) insert(k []int64, v uint64) (domains, bool) {
	if !d.ready {
		panic(fmt.Errorf("run is not ready for inplace insert %d (%x)", len(k), v))
	}

	if r, ok := d.data.insert(k, v); ok {
		d.data = r
		return d, true
	}

	return d, false
}

func (d domains) inplaceInsert(k []int64, v uint64) domains {
	if !d.ready {
		panic(fmt.Errorf("run is not ready for inplace insert %d (%x)", len(k), v))
	}

	d.data = d.data.inplaceInsert(k, v)
	return d
}

func (d domains) append(k []int64, v uint64) domains {
	r := d.data
	if d.ready {
		r = r.clone()
	}

	d.data = r.append(k, v)
	d.ready = false
	return d
}

func (d domains) normalize(log normalizeLoggers) domains {
	if d.ready {
		return d
	}

	d.data = d.data.normalize(log)
	d.ready = true
	return d
}

func (d domains) get(k []int64) uint64 {
	if !d.ready {
		panic(fmt.Errorf("run is not ready for get %d", len(k)))
	}

	v, start, end := d.data.get(k)
	if end != 0 {
		v, err := d.getter.get(k, start, end)
		if err != nil {
			panic(fmt.Errorf("can't read data for range [%d, %d) from %s: %s", start, end, d.getter.path, err))
		}

		return v
	}

	return v
}

func (d domains) size() int {
	return d.data.size()
}

func (d domains) size90() int {
	return d.data.size90()
}

func (d domains) flush(fLog flushLoggers, rLog func(size, from, to, reqs, queue int), nLog normalizeLoggers) (domains, *Getter) {
	var getter *Getter

	if d.dir == nil {
		return d, getter
	}

	n := len(d.data.v)
	if n > 0 {
		if len(d.data.k) < n {
			panic(fmt.Errorf("corrupted run (k: %d, v: %d) on flush", len(d.data.k), n))
		}

		m := len(d.data.k) / n

		if fLog.before != nil {
			fLog.before(m, n, n/10)
		}

		d = d.normalize(nLog)
		tmp, g, err := d.merge(rLog)
		if err != nil {
			panic(err)
		}

		d = tmp.drop90()
		getter = g

		if fLog.after != nil {
			fLog.after(m, n, len(d.data.v))
		}
	}

	return d, getter
}

func (d domains) merge(log func(size, from, to, reqs, queue int)) (domains, *Getter, error) {
	var getter *Getter

	n := len(d.data.v)
	if n > 0 {
		if len(d.data.k) < n {
			panic(fmt.Errorf("corrupted run (k: %d, v: %d) on write all", len(d.data.k), n))
		}

		m := len(d.data.k) / n
		dst, err := ioutil.TempFile(*d.dir, fmt.Sprintf("%03d.", m))
		if err != nil {
			return d, getter, err
		}

		path := dst.Name()
		defer func() {
			if err != nil {
				dst.Close()
				os.Remove(path)
			}
		}()

		b := bufio.NewWriter(dst)
		w := d.makeWriterRun(b)

		r, c, err := d.openStorage()
		if err != nil {
			return d, getter, err
		}
		if c != nil {
			defer c.Close()
		}

		filter, data, err := d.data.merge(r, w)
		if err != nil {
			return d, getter, err
		}

		err = b.Flush()
		if err != nil {
			return d, getter, err
		}

		err = dst.Close()
		if err != nil {
			return d, getter, err
		}

		g := newGetter(path, m, w.blks, w.rem, filter, log)
		err = g.start()
		if err != nil {
			return d, getter, err
		}

		getter = d.getter

		d.getter = g
		d.data = data
		d.blks = w.blks
		d.rem = w.rem
	}

	return d, getter, nil
}

func (d domains) openStorage() (*rRun, io.Closer, error) {
	if d.getter == nil {
		return d.makeReaderRun(nil), nil, nil
	}

	f, err := os.Open(d.getter.path)
	if err != nil {
		return nil, nil, err
	}

	return d.makeReaderRun(bufio.NewReader(f)), f, nil
}

func (d domains) drop90() domains {
	d.data = d.data.drop90()
	return d
}
