package sdntable64

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
)

type domains struct {
	path  *string
	ready bool
	data  run
	blks  int
	rem   int
}

func makeDomains() domains {
	return domains{
		ready: true,
	}
}

func makeDomainsWithPath(path string, n int) domains {
	d := makeDomains()

	path = filepath.Join(path, fmt.Sprintf("%03d", n))
	d.path = &path

	return d
}

func (d domains) inplaceInsert(k []int64, v uint64) domains {
	if !d.ready {
		panic(fmt.Errorf("run is not ready for inplace insert %d (%x)", len(k), v))
	}

	d.data = d.data.inplaceInsert(k, v)
	return d
}

func (d domains) append(k []int64, v uint64) domains {
	d.data = d.data.append(k, v)
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
		tmp, err := d.readFromDisk(start, end, read)
		if err != nil {
			panic(fmt.Errorf("can't read data for range [%d, %d) from %s: %s", start, end, *d.path, err))
		}

		v, _, _ = tmp.get(k)
	}

	return v
}

func (d domains) size() int {
	return d.data.size()
}

func (d domains) size90() int {
	return d.data.size90()
}

func (d domains) flush(fLog flushLoggers, nLog normalizeLoggers) domains {
	if d.path == nil {
		return d
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
			tmp, err := d.merge()
			if err != nil {
				panic(err)
			}

			d = tmp
		} else {
			tmp, err := d.writeAll()
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

	return d
}

func (d domains) readFromDisk(start, end uint32, log func(size, from, to int)) (run, error) {
	out := run{
		k: []int64{},
		v: []uint64{},
	}

	f, err := os.Open(*d.path)
	if err != nil {
		return out, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	n := len(d.data.v)
	m := len(d.data.k) / n
	blk := blocks[m-1]

	buf := makeRunForBlock(m)

	cur := int(start / uint32(blk))
	last := d.blks
	if end < math.MaxUint32 {
		last = int(end / uint32(blk))
	}

	if log != nil {
		log(m, cur, last)
	}

	if cur > 0 {
		if _, err := f.Seek(int64(cur)*int64(blk), 0); err != nil {
			return out, err
		}
	}

	for cur < d.blks && cur <= last {
		if err := buf.read(r); err != nil {
			return out, err
		}

		out.k = append(out.k, buf.k...)
		out.v = append(out.v, buf.v...)

		cur++
	}

	if cur < last && cur >= d.blks {
		buf = buf.truncate(d.rem)
		if err := buf.read(r); err != nil {
			return out, err
		}

		out.k = append(out.k, buf.k...)
		out.v = append(out.v, buf.v...)
	}

	return out, nil
}

func (d domains) merge() (domains, error) {
	dst, err := ioutil.TempFile("", "*")
	if err != nil {
		return d, err
	}
	defer func() {
		dst.Close()
		if err == nil {
			os.Rename(dst.Name(), *d.path)
		} else {
			os.Remove(dst.Name())
		}
	}()

	b := bufio.NewWriter(dst)
	defer b.Flush()

	w := d.makeWriterRun(b)

	src, err := os.Open(*d.path)
	if err != nil {
		return d, err
	}
	defer src.Close()

	r, err := d.data.merge(d.makeReaderRun(bufio.NewReader(src)), w)
	if err != nil {
		return d, err
	}

	d.data = r
	d.blks = w.blks
	d.rem = w.rem

	return d, nil
}

func (d domains) writeAll() (domains, error) {
	f, err := ioutil.TempFile("", "*")
	if err != nil {
		return d, err
	}
	defer func() {
		f.Close()
		if err == nil {
			os.Rename(f.Name(), *d.path)
		} else {
			os.Remove(f.Name())
		}
	}()

	w := bufio.NewWriter(f)
	defer w.Flush()

	r, blks, rem, err := d.data.writeAll(w)
	if err != nil {
		return d, err
	}

	d.data = r
	d.blks = blks
	d.rem = rem

	return d, nil
}

func (d domains) drop90() domains {
	d.data = d.data.drop90()
	return d
}
