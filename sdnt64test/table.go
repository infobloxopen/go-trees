package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/infobloxopen/go-trees/sdntable64"
	"github.com/infobloxopen/go-trees/udomain"
)

type table struct {
	t *sdntable64.Table64
}

func newTable(in []pair, dir string, reqsStat *[]reqStat) (mapper64, error) {
	opts := []sdntable64.Option{
		sdntable64.WithNormalizeLoggers(func(i, j int) {
			log.Printf("normalizing %d subarray with %d domains", i, j)
		}, func(i, j int) {
			log.Printf("sorting %d subarray with %d domains", i, j)
		}, func(i, j int) {
			log.Printf("normalized %d subarray with %d domains", i, j)
		}),
		sdntable64.WithFlushLoggers(func(i, j, k int) {
			log.Printf("flushing %d subarray (%d -> %d)", i, j, k)
		}, func(i, j, k int) {
			log.Printf("flushed %d subarray (%d -> %d)", i, j, k)
		}),
	}

	if reqsStat != nil && *reqsStat != nil {
		opts = append(opts,
			sdntable64.WithReadLogger(func(size, from, to, reqs, queue int) {
				*reqsStat = append(*reqsStat, reqStat{
					reqs:  reqs,
					queue: queue,
				})
			}),
		)
	}

	if len(dir) > 0 {
		path := filepath.Join(os.TempDir(), dir)
		if err := os.RemoveAll(path); err != nil {
			log.Fatalf("can't cleanup directory for table storage %q: %s", path, err)
		}

		if err := os.MkdirAll(path, 0777); err != nil {
			log.Fatalf("can't create directory for table storage %q: %s", path, err)
		}

		opts = append(opts, sdntable64.WithPath(path))
	}

	t := sdntable64.NewTable64(opts...)
	for i, p := range in {
		d, err := domain.MakeNameFromString(p.k)
		if err != nil {
			return nil, fmt.Errorf("can't convert %q at %d to domain name: %s", p.k, i+1, err)
		}

		tmp, g := t.Append(d, p.v)
		for _, g := range g {
			if err := g.Stop(); err != nil {
				return nil, err
			}
		}
		t = tmp

		if (i+1)%1000000 == 0 {
			log.Printf("inserted %d domains", i+1)
		}
	}

	if len(in)%1000000 != 0 {
		log.Printf("inserted %d domains", len(in))
	}

	t = t.Normalize()
	log.Printf("Table size: %s", makeSize(uint64(t.Size())))

	return &table{
		t: t.Normalize(),
	}, nil
}

func (t *table) Map(k string) uint64 {
	d, err := domain.MakeNameFromString(k)
	if err != nil {
		log.Fatalf("can't convert %q to domain name: %s", k, err)
	}

	return t.t.Get(d)
}
