package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

type config struct {
	path    string
	data    string
	pause   time.Duration
	reqs    int
	workers int
}

var conf config

const (
	mapDataStruct      = "map"
	tableDataStruct    = "table"
	memTableDataStruct = "mem-table"
)

var dataStructs = map[string]struct{}{
	mapDataStruct:      {},
	tableDataStruct:    {},
	memTableDataStruct: {},
}

func init() {
	flag.StringVar(&conf.path, "d", "domains.lst", "thread feed categories by domains")
	flag.StringVar(&conf.data, "s", mapDataStruct,
		fmt.Sprintf("data structure to test "+
			"%q - map[string]uint64, "+
			"%q - sdntable64, "+
			"%q - in-memory sdntable64",
			mapDataStruct, tableDataStruct, memTableDataStruct))
	flag.DurationVar(&conf.pause, "p", 0, "pause before exit")
	flag.IntVar(&conf.reqs, "n", -1, "number of domains to request, <0 - request all domains")
	flag.IntVar(&conf.workers, "w", 0, "number of workers, <1 - make requests synchronously")
	flag.Parse()

	if _, ok := dataStructs[conf.data]; !ok {
		log.Fatalf("expected %q, %q or %q as data structure but got %q",
			mapDataStruct, tableDataStruct, memTableDataStruct, conf.data)
	}
}
