package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func load(name string) map[interface{}]interface{} {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := yaml.NewDecoder(f)

	out := make(map[interface{}]interface{})
	if err := d.Decode(out); err != nil {
		log.Fatalf("can't decode data from %q: %s", name, err)
	}

	return out
}
