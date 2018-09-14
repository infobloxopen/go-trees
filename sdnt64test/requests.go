package main

import (
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func run(in []pair, miss []string, m mapper64) {
	count := len(in)
	if conf.reqs > 0 {
		count = conf.reqs
	}
	report := makeReport(count)

	var (
		idx  uint64
		dIdx uint64
		mIdx uint64
	)

	start := time.Now()

	if conf.workers > 0 {
		wg := new(sync.WaitGroup)
		wg.Add(conf.workers)

		for i := 0; i < conf.workers; i++ {
			go func(wg *sync.WaitGroup, pIdx, pDIdx, pMIdx *uint64) {
				defer wg.Done()

				for {
					i := int(atomic.AddUint64(pIdx, 1) - 1)
					if i >= count {
						return
					}

					if len(miss) <= 0 || rand.Float64()*100 >= conf.missPart {
						j := 0
						if conf.rand {
							j = rand.Intn(len(in))
						} else {
							j = int(atomic.AddUint64(pDIdx, 1) - 1)
						}
						p := in[j%len(in)]

						report[i].start = time.Now()
						if u := m.Map(p.k); u != p.v {
							log.Fatalf("invalid result for %q (%d): %x != %x", p.k, i+1, u, p.v)
						}
						report[i].end = time.Now()
					} else {
						j := 0
						if conf.rand {
							j = rand.Intn(len(miss))
						} else {
							j = int(atomic.AddUint64(pMIdx, 1) - 1)
						}
						k := miss[j%len(miss)]

						report[i].start = time.Now()
						m.Map(k)
						report[i].end = time.Now()
					}
				}
			}(wg, &idx, &dIdx, &mIdx)
		}

		wg.Wait()
	} else {
		for i := 0; i < count; i++ {
			if len(miss) <= 0 || rand.Float64()*100 >= conf.missPart {
				j := 0
				if conf.rand {
					j = rand.Intn(len(in))
				} else {
					dIdx++
					j = int(dIdx - 1)
				}
				p := in[j%len(in)]

				report[i].start = time.Now()
				if u := m.Map(p.k); u != p.v {
					log.Fatalf("invalid result for %q (%d): %x != %x", p.k, i+1, u, p.v)
				}
				report[i].end = time.Now()
			} else {
				j := 0
				if conf.rand {
					j = rand.Intn(len(miss))
				} else {
					mIdx++
					j = int(mIdx - 1)
				}
				k := miss[j%len(miss)]

				report[i].start = time.Now()
				m.Map(k)
				report[i].end = time.Now()
			}
		}
	}

	dur := time.Now().Sub(start)
	dt := dur.Nanoseconds() / int64(count)
	log.Printf("Duration: %s (%d ns/op; %d op/s)", dur, dt, time.Second.Nanoseconds()/dt)

	report.dump(start)
}
