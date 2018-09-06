package main

import "flag"

type config struct {
	input  string
	output string
	before uint
	after  uint
	step   uint
}

var conf config

func init() {
	flag.StringVar(&conf.input, "d", "domains.lst", "list of domain names")
	flag.StringVar(&conf.output, "m", "missing.lst", "list of generated missing domain names")
	flag.UintVar(&conf.before, "b", 50, "number of domain names to generate before first domain")
	flag.UintVar(&conf.after, "a", 50, "number of domain names to generate after last domain")
	flag.UintVar(&conf.step, "s", 1000, "step to find pair to generate domain in between")

	flag.Parse()
}
