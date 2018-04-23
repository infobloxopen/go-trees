package main

import (
	"flag"
	"log"
)

type config struct {
	template string
	data     string
}

var conf config

func parse() {
	flag.StringVar(&conf.template, "t", "", "path to template (required)")
	flag.StringVar(&conf.data, "d", "", "path to data")

	flag.Parse()

	if len(conf.template) <= 0 {
		log.Fatal("No path to template - nothing to execute")
	}
}
