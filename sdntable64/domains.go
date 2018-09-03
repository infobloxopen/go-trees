package sdntable64

import (
	"bufio"
	"fmt"
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

func (d domains) get(k []int64, read func(size, from, to int)) uint64 {
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

func (d domains) flush(fLog flushLoggers, rLog func(size, from, to int), nLog normalizeLoggers) (domains, *Getter) {
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
		if len(d.data.i) > 0 {
			tmp, g, err := d.merge(rLog)
			if err != nil {
				panic(err)
			}

			d = tmp
			getter = g
		} else {
			tmp, err := d.writeAll(rLog)
			if err != nil {
				panic(err)
			}

			d = tmp
		}

		d = d.drop90()

		if fLog.after != nil {
			fLog.after(m, n, len(d.data.v))
		}
	}

	return d, getter
}

func (d domains) merge(log func(size, from, to int)) (domains, *Getter, error) {
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

		src, err := os.Open(d.getter.path)
		if err != nil {
			return d, getter, err
		}
		defer src.Close()

		r, err := d.data.merge(d.makeReaderRun(bufio.NewReader(src)), w)
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

		g := newGetter(path, m, w.blks, w.rem, log)
		err = g.start()
		if err != nil {
			return d, getter, err
		}

		getter = d.getter

		d.getter = g
		d.data = r
		d.blks = w.blks
		d.rem = w.rem
	}

	return d, getter, nil
}

func (d domains) writeAll(log func(size, from, to int)) (domains, error) {
	n := len(d.data.v)
	if n > 0 {
		if len(d.data.k) < n {
			panic(fmt.Errorf("corrupted run (k: %d, v: %d) on write all", len(d.data.k), n))
		}

		m := len(d.data.k) / n
		f, err := ioutil.TempFile(*d.dir, fmt.Sprintf("%03d.", m))
		if err != nil {
			return d, err
		}
		path := f.Name()
		defer func() {
			if err != nil {
				f.Close()
				os.Remove(path)
			}
		}()

		w := bufio.NewWriter(f)

		r, blks, rem, err := d.data.writeAll(w)
		if err != nil {
			return d, err
		}

		err = w.Flush()
		if err != nil {
			return d, err
		}

		err = f.Close()
		if err != nil {
			return d, err
		}

		g := newGetter(path, m, blks, rem, log)
		err = g.start()
		if err != nil {
			return d, err
		}

		d.getter = g
		d.data = r
		d.blks = blks
		d.rem = rem
	}

	return d, nil
}

func (d domains) drop90() domains {
	d.data = d.data.drop90()
	return d
}
