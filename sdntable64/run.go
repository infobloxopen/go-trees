package sdntable64

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type run struct {
	k []int64
	v []uint64
	i []uint32
}

func makeRun(k []int64, v uint64) run {
	keys := make([]int64, len(k))
	copy(keys, k)

	return run{
		k: keys,
		v: []uint64{v},
	}
}

func makeRunForBlock(n int) run {
	m := blocks[n-1]
	return run{
		k: make([]int64, m*n),
		v: make([]uint64, m),
	}
}

func (r run) inplaceInsert(k []int64, v uint64) run {
	if len(r.v) > 0 {
		if len(r.k) < len(r.v) {
			panic(fmt.Errorf("corrupted run (k: %d, v: %d) on inplace insert %d (%x)", len(r.k), len(r.v), len(k), v))
		}

		if len(r.k)/len(r.v) != len(k) {
			panic(fmt.Errorf("invalid key on inplace insert %d (%x) for run (k: %d, v: %d, n: %d)",
				len(k), v, len(r.k), len(r.v), len(r.k)/len(r.v)),
			)
		}

		left := 0
		right := len(r.v)
		for {
			m := (left + right) / 2
			b := m * len(k)
			d := compare(k, r.k[b:b+len(k)])
			if d == 0 {
				r.v[m] = v
				return r
			}

			if d < 0 {
				right = m
				if left == right {
					break
				}
			} else {
				if left == m {
					left = right
					break
				}

				left = m
			}
		}

		r.k = append(r.k, k...)
		r.v = append(r.v, v)
		if len(r.i) > 0 {
			r.i = append(r.i, math.MaxUint32)
		}

		if right < len(r.v)-1 {
			i := len(k) * right
			copy(r.k[i+len(k):], r.k[i:])
			copy(r.k[i:], k)

			copy(r.v[right+1:], r.v[right:])
			r.v[right] = v

			if len(r.i) > 0 {
				copy(r.i[right+1:], r.i[right:])
				r.i[right] = math.MaxUint32
			}
		}

		return r
	}

	return makeRun(k, v)
}

func (r run) append(k []int64, v uint64) run {
	if len(r.v) > 0 {
		if len(r.k) < len(r.v) {
			panic(fmt.Errorf("corrupted run (k: %d, v: %d) on append %d (%x)", len(r.k), len(r.v), len(k), v))
		}

		if len(r.k)/len(r.v) != len(k) {
			panic(fmt.Errorf("invalid key on append %d (%x) for run (k: %d, v: %d, n: %d)",
				len(k), v, len(r.k), len(r.v), len(r.k)/len(r.v)),
			)
		}

		r.k = append(r.k, k...)
		r.v = append(r.v, v)
		if len(r.i) > 0 {
			r.i = append(r.i, math.MaxUint32)
		}
		return r
	}

	return makeRun(k, v)
}

func (r run) normalize(log normalizeLoggers) run {
	n := len(r.v)
	if n > 0 {
		if len(r.k) < n {
			panic(fmt.Errorf("corrupted run (k: %d, v: %d) on normalize", len(r.k), len(r.v)))
		}

		m := len(r.k) / n

		if n > 1 {
			if log.before != nil {
				log.before(m, n)
			}

			ir := r.makeIndexedRun()

			if log.sort != nil {
				log.sort(m, n)
			}
			ir.sort()

			r = ir.makeRun()

			if log.after != nil {
				log.after(m, len(r.v))
			}
		}
	}

	return r
}

func (r run) get(k []int64) (uint64, uint32, uint32) {
	n := len(r.v)
	if n > 0 {
		if len(r.k) < n {
			panic(fmt.Errorf("corrupted run (k: %d, v: %d) on get %d", len(r.k), n, len(k)))
		}

		m := len(r.k) / n
		if m != len(k) {
			panic(fmt.Errorf("invalid key on get %d for run (k: %d, v: %d, n: %d)", len(k), len(r.k), n, m))
		}

		left := 0
		right := len(r.v)
		for {
			m := (left + right) / 2
			b := m * len(k)
			d := compare(k, r.k[b:b+len(k)])
			if d == 0 {
				return r.v[m], 0, 0
			}

			if left == m {
				if d < 0 {
					right = m
				}
				break
			}

			if d < 0 {
				right = m
			} else {
				left = m
			}
		}

		if len(r.i) > 0 {
			start, end := r.getDiskRange(right)
			return 0, start, end
		}
	}

	return 0, 0, 0
}

func (r run) getDiskRange(right int) (uint32, uint32) {
	start := uint32(0)
	end := uint32(math.MaxUint32)

	left := right
	if left > 0 {
		left--
		for left >= 0 && (left >= len(r.i) || r.i[left] >= math.MaxUint32) {
			left--
		}

		if left >= 0 {
			start = r.i[left]
		}
	}

	for right < len(r.i) && r.i[right] >= math.MaxUint32 {
		right++
	}

	if right < len(r.i) {
		end = r.i[right]
	}

	return start, end
}

func (r run) size() int {
	return (len(r.k) + len(r.v)) * 8
}

func (r run) size90() int {
	if len(r.v) > 0 {
		if len(r.k) < len(r.v) {
			panic(fmt.Errorf("corrupted run (k: %d, v: %d) on geting size", len(r.k), len(r.v)))
		}

		n := len(r.v)
		m := len(r.k) / n

		return n * 90 / 100 * 8 * (m + 1)
	}

	return 0
}

func (r run) truncate(n int) run {
	if n >= 0 && len(r.v) > n {
		if len(r.k) < len(r.v) {
			panic(fmt.Errorf("corrupted run (k: %d, v: %d) on truncate", len(r.k), len(r.v)))
		}

		m := len(r.k) / len(r.v)

		r.k = r.k[:m*n]
		r.v = r.v[:n]
	}

	return r
}

func (r run) read(f io.Reader) error {
	if err := binary.Read(f, binary.LittleEndian, r.k); err != nil {
		return err
	}

	if err := binary.Read(f, binary.LittleEndian, r.v); err != nil {
		return err
	}

	return nil
}

func (r run) write(f io.Writer) error {
	if err := binary.Write(f, binary.LittleEndian, r.k); err != nil {
		return err
	}

	if err := binary.Write(f, binary.LittleEndian, r.v); err != nil {
		return err
	}

	return nil
}

func (r run) writeAll(f io.Writer) (run, int, int, error) {
	n := len(r.v)
	m := len(r.k) / n

	blk := blocks[m-1]
	blks := 0

	k := r.k
	v := r.v

	idx := make([]uint32, n)
	for i := range idx {
		if i >= math.MaxUint32 {
			return r, blks, n, fmt.Errorf("run %d overflow", m)
		}
		idx[i] = uint32(i)
	}

	for n >= blk {
		tmp := run{
			k: k[:m*blk],
			v: v[:blk],
		}
		if err := tmp.write(f); err != nil {
			return r, blks, n, err
		}

		k = k[m*blk:]
		v = v[blk:]

		n -= blk
		blks++
	}

	if n > 0 {
		if err := r.write(f); err != nil {
			return r, blks, n, err
		}
	}

	r.i = idx
	return r, blks, n, nil
}

func (r run) merge(in *rRun, out *wRun) (run, error) {
	n := len(r.v)
	m := len(r.k) / n
	idx := make([]uint32, n)

	mK := r.k
	i := 0

	sK, sV, err := in.next()
	if err != nil {
		return r, err
	}

	var c uint32
	for len(mK) > 0 || len(sK) > 0 {
		var d int64 = 1
		if len(mK) > 0 && len(sK) > 0 {
			d = compare(mK[:m], sK)
		} else if len(mK) > 0 {
			d = -1
		}

		if d <= 0 {
			if err := out.put(mK[:m], r.v[i]); err != nil {
				return r, err
			}

			mK = mK[m:]
			idx[i] = c
			i++
		} else {
			if err := out.put(sK, sV); err != nil {
				return r, err
			}
		}

		if d >= 0 {
			sK, sV, err = in.next()
			if err != nil {
				return r, err
			}
		}

		c++
		if c >= math.MaxUint32 {
			return r, fmt.Errorf("run %d overflow", m)
		}
	}

	if err := out.flush(); err != nil {
		return r, err
	}

	r.i = idx
	return r, nil
}

func (r run) drop90() run {
	n := len(r.v)
	if n > 0 {
		if len(r.k) < n {
			panic(fmt.Errorf("corrupted run (k: %d, v: %d) on dropping 90%% entries", len(r.k), n))
		}

		m := len(r.k) / n

		n /= 10
		keys := make([]int64, n*m)
		values := make([]uint64, n)
		var idx []uint32
		if len(r.i) > 0 {
			idx = make([]uint32, n)
		}

		j := 0
		for i := range values {
			a := j * m
			copy(keys[i*m:], r.k[a:a+m])
			values[i] = r.v[j]
			if idx != nil {
				idx[i] = r.i[j]
			}

			j += 10
		}

		r = run{
			k: keys,
			v: values,
			i: idx,
		}
	}

	return r
}

var blocks = []int{
	6400,
	4266,
	3200,
	2560,
	2133,
	1828,
	1600,
	1422,
	1280,
	1163,
	1066,
	984,
	914,
	853,
	800,
	752,
	711,
	673,
	640,
	609,
	581,
	556,
	533,
	512,
	492,
	474,
	457,
	441,
	426,
	412,
	400,
	387,
	376,
	365,
	355,
	345,
	336,
	328,
	320,
	312,
	304,
	297,
	290,
	284,
	278,
	272,
	266,
	261,
	256,
	250,
	246,
	241,
	237,
	232,
	228,
	224,
	220,
	216,
	213,
	209,
	206,
	203,
	200,
	196,
	193,
	191,
	188,
	185,
	182,
	180,
	177,
	175,
	172,
	170,
	168,
	166,
	164,
	162,
	160,
	158,
	156,
	154,
	152,
	150,
	148,
	147,
	145,
	143,
	142,
	140,
	139,
	137,
	136,
	134,
	133,
	131,
	130,
	129,
	128,
	126,
	125,
	124,
	123,
	121,
	120,
	119,
	118,
	117,
	116,
	115,
	114,
	113,
	112,
	111,
	110,
	109,
	108,
	107,
	106,
	105,
	104,
	104,
	103,
	102,
	101,
	100,
	100,
}
