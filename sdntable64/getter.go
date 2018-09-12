package sdntable64

import (
	"encoding/binary"
	"io"
	"math"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/golang-lru"
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
	cache  *lru.Cache
	log    func(size, from, to, reqs, queue int)
}

func newGetter(path string, m, blks, rem int, filter *bloom.BloomFilter, log func(size, from, to, reqs, queue int)) *Getter {
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

	buf := run{
		k: make([]int64, g.m*g.blk),
		v: make([]uint64, g.blk),
	}

	queue := map[uint64][]getRequest{}

	rPool := sync.Pool{
		New: func() interface{} {
			return run{
				k: make([]int64, 0, g.m*g.blk),
				v: make([]uint64, g.blk),
			}
		},
	}
	cache, err := lru.NewWithEvict(500, func(k, v interface{}) {
		rPool.Put(v)
	})
	if err != nil {
		return err
	}
	g.cache = cache

	g.wg.Add(1)
	go func(wg *sync.WaitGroup, ch chan getRequest, done chan struct{}) {
		iwg := new(sync.WaitGroup)
		defer func() {
			iwg.Wait()
			wg.Done()
		}()

		t := time.NewTicker(100 * time.Microsecond)

		for {
			select {
			case <-done:
				t.Stop()
				return

			case <-t.C:
				var (
					maxQ []getRequest
					maxK uint64
				)

				for k, q := range queue {
					if len(q) > len(maxQ) {
						maxQ = q
						maxK = k
					}
				}

				if len(maxQ) > 0 {
					iwg.Wait()

					delete(queue, maxK)

					iwg.Add(1)
					go func(wg *sync.WaitGroup, k uint64, q []getRequest, lenQ int) {
						defer wg.Done()

						r := q[0]

						if g.log != nil {
							g.log(len(r.k), int(r.start), int(r.end), len(q), lenQ)
						}

						tmp, err := g.read(f, r.start, r.end, rPool.Get().(run), buf, nil)
						if err != nil {
							for _, r := range q {
								r.ch <- getResponse{err: err}
							}

							return
						}

						g.cache.Add(k, tmp)

						for _, r := range q {
							u, _, _ := tmp.get(r.k)
							r.ch <- getResponse{value: u}
						}

					}(iwg, maxK, maxQ, len(queue))
				}

			case r := <-ch:
				k := (uint64(r.start) << 32) | uint64(r.end)
				q, ok := queue[k]
				if !ok {
					q = []getRequest{r}
				} else {
					q = append(q, r)
				}
				queue[k] = q
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

	idx := (uint64(start) << 32) | uint64(end)
	if tmp, ok := g.cache.Get(idx); ok {
		if u, _, _ := tmp.(run).get(k); u != 0 {
			return u, nil
		}
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
