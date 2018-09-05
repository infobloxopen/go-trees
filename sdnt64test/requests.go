package main

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

func run(in []pair, m mapper64) {
	count := len(in)
	if conf.reqs > 0 && conf.reqs < count {
		count = conf.reqs
	}
	report := makeReport(count)

	var idx uint64

	start := time.Now()

	wg := new(sync.WaitGroup)
	wg.Add(conf.workers)

	for i := 0; i < conf.workers; i++ {
		go func(wg *sync.WaitGroup, pIdx *uint64) {
			defer wg.Done()

			for {
				i := int(atomic.AddUint64(pIdx, 1) - 1)
				if i >= count {
					return
				}

				p := in[i%len(in)]

				report[i].start = time.Now()
				if u := m.Map(p.k); u != p.v {
					log.Fatalf("invalid result for %q (%d): %x != %x", p.k, i+1, u, p.v)
				}
				report[i].end = time.Now()
			}
		}(wg, &idx)
	}

	wg.Wait()

	dur := time.Now().Sub(start)
	dt := dur.Nanoseconds() / int64(count)
	log.Printf("Duration: %s (%d ns/op; %d op/s)", dur, dt, time.Second.Nanoseconds()/dt)

	report.dump(start)
}
