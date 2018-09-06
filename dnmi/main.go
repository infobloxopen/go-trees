package main

import "log"

func main() {
	ss, err := load(conf.input)
	if err != nil {
		log.Fatalf("failed to load domains from %s: %s", conf.input, err)
	}

	log.Printf("loaded %d domains from %s", len(ss), conf.input)

	ds, err := convert(ss)
	if err != nil {
		log.Fatalf("failed to convert domains: %s", err)
	}

	ms, err := generate(ds, int(conf.before), int(conf.step), int(conf.after))
	if err != nil {
		log.Fatalf("failed to generate missing domains: %s", err)
	}

	log.Printf("generated totally %d domains", len(ms))

	if err := save(conf.output, ms); err != nil {
		log.Fatalf("failed to save domains: %s", err)
	}
}
