package sdntable64

import (
	"encoding/binary"
	"io"
	"math"
	"os"
	"sync"
	"syscall"

	"github.com/willf/bloom"
)

type Getter struct {
	path   string
	f      io.Closer
	filter *bloom.BloomFilter
	m      int
	blk    int
	blks   int
	rem    int
	bSize  int64
	ch     chan getRequest
	done   chan struct{}
	wg     *sync.WaitGroup
	pool   chan chan getResponse
	log    func(size, from, to int)
}

func newGetter(path string, m, blks, rem int, filter *bloom.BloomFilter, log func(size, from, to int)) *Getter {
	blk := blocks[m-1]
	pool := make(chan chan getResponse, 128)
	for len(pool) < cap(pool) {
		pool <- make(chan getResponse)
	}

	return &Getter{
		path:   path,
		filter: filter,
		m:      m,
		blk:    blk,
		blks:   blks,
		rem:    rem,
		bSize:  int64(blk) * 8 * (int64(m) + 1),
		ch:     make(chan getRequest),
		done:   make(chan struct{}),
		wg:     new(sync.WaitGroup),
		pool:   pool,
		log:    log,
	}
}

func (g *Getter) start() error {
	f, err := os.Open(g.path)
	if err != nil {
		return err
	}

	g.f = f

	tmp := run{
		k: []int64{},
		v: []uint64{},
	}

	buf := run{
		k: make([]int64, g.m*g.blk),
		v: make([]uint64, g.blk),
	}

	g.wg.Add(1)
	go func(wg *sync.WaitGroup, ch chan getRequest, done chan struct{}) {
		defer wg.Done()

		for {
			select {
			case <-done:
				return

			case r := <-ch:
				tmp, err := g.read(f, r.start, r.end, tmp, buf, g.log)
				if err != nil {
					r.ch <- getResponse{err: err}
					break
				}

				u, _, _ := tmp.get(r.k)
				r.ch <- getResponse{value: u}
			}
		}
	}(g.wg, g.ch, g.done)

	return nil
}

func (g *Getter) read(f io.ReadSeeker, start, end uint32, tmp, buf run, log func(size, from, to int)) (run, error) {
	cur := int(start / uint32(g.blk))
	last := g.blks
	if end < math.MaxUint32 {
		last = int(end / uint32(g.blk))
	}

	if log != nil {
		log(g.m, cur, last)
	}

	if _, err := f.Seek(int64(cur)*g.bSize, io.SeekStart); err != nil {
		return tmp, err
	}

	tmp.k = tmp.k[:0]
	tmp.v = tmp.v[:0]

	for cur < g.blks && cur <= last {
		if err := buf.read(f); err != nil {
			return tmp, err
		}

		tmp.k = append(tmp.k, buf.k...)
		tmp.v = append(tmp.v, buf.v...)

		cur++
	}

	if cur <= last && cur >= g.blks {
		buf = buf.truncate(g.rem)
		if err := buf.read(f); err != nil {
			return tmp, err
		}

		tmp.k = append(tmp.k, buf.k...)
		tmp.v = append(tmp.v, buf.v...)
	}

	return tmp, nil
}

func (g *Getter) Stop() error {
	close(g.done)
	g.wg.Wait()

	if g.f != nil {
		if err := g.f.Close(); err != nil {
			return err
		}
	}

	if err := os.Remove(g.path); err != nil {
		if err, ok := err.(*os.PathError); ok {
			if err.Err != syscall.ENOENT {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

type getRequest struct {
	k     []int64
	start uint32
	end   uint32
	ch    chan getResponse
}

func (g *Getter) get(k []int64, start, end uint32) (uint64, error) {
	b := make([]byte, 8*len(k))
	for i, n := range k {
		binary.LittleEndian.PutUint64(b[8*i:], uint64(n))
	}
	if g.filter != nil && !g.filter.Test(b) {
		return 0, nil
	}

	ch := <-g.pool
	defer func() {
		g.pool <- ch
	}()

	g.ch <- getRequest{
		k:     k,
		start: start,
		end:   end,
		ch:    ch,
	}

	r := <-ch
	return r.value, r.err
}

type getResponse struct {
	value uint64
	err   error
}
