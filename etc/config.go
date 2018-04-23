package main

import (
	"flag"
	"log"
)

type config struct {
	template string
}

var conf config

func init() {
	flag.StringVar(&conf.template, "t", "", "path to template (required)")

	flag.Parse()

	if len(conf.template) <= 0 {
		log.Fatal("No path to template - nothing to execute")
	}
}
