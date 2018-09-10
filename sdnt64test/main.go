package main

import (
	"log"
	"time"
)

func main() {
	total := 0
	var pairs []pair

	if err := load(conf.path, func(n int) error {
		log.Printf("loading %d domains from %q", n, conf.path)
		pairs = make([]pair, 0, n)
		return nil
	}, func(k string, v uint64) error {
		total += len(k)

		pairs = append(pairs, pair{
			k: k,
			v: v,
		})

		if len(pairs)%1000000 == 0 {
			log.Printf("loaded %d domains", len(pairs))
		}

		return nil
	}); err != nil {
		log.Fatalf("can't read feed from %q: %s", conf.path, err)
	}

	log.Printf("got %d domains from %q (avg. domain length %.02f)",
		len(pairs), conf.path, float64(total)/float64(len(pairs)))

	printAlloc("array")

	var miss []string
	if len(conf.miss) > 0 && conf.missPart > 0 {
		log.Printf("loading missing domains from %q", conf.miss)
		s, err := loadMissing(conf.miss)
		if err != nil {
			log.Fatalf("can't read domains from %q: %s", conf.miss, err)
		}

		miss = s
		log.Printf("loaded %d domains from %q", len(miss), conf.miss)

		count := len(miss)
		if conf.reqs > 0 && conf.reqs < count {
			count = conf.reqs
		}
		log.Printf("going to make %.02f requests from missing domains list", float64(count)*conf.missPart/100)
	}

	printAlloc("missing")

	var (
		m   mapper64
		err error
	)

	switch conf.data {
	default:
		m, err = newTable(pairs, "sdntable64")
		if err != nil {
			log.Fatalf("can't fill table: %s", err)
		}

		printAlloc(conf.data)
		log.Printf("Table: %p", m)

	case memTableDataStruct:
		m, err = newTable(pairs, "")
		if err != nil {
			log.Fatalf("can't fill table: %s", err)
		}

		printAlloc(conf.data)
		log.Printf("Table: %p", m)

	case mapDataStruct:
		m = newMap(pairs)

		printAlloc(conf.data)
		log.Printf("Map: %p", m)
	}

	run(pairs, miss, m)

	if conf.pause > 0 {
		log.Printf("paused for %s before exit", conf.pause)
		time.Sleep(conf.pause)
	}
}
