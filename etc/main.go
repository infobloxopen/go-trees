// ETC (Execute Template Catalog) executes all templates recursively in a given directory.
package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	execute(conf.template, makePrefix(conf.template))
}

func execute(name string, prefix []string) {
	absName, err := filepath.Abs(name)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("executing %q", getRelName(absName, prefix))

	fi, err := os.Stat(absName)
	if err != nil {
		log.Fatalf("\t%s", err)
	}

	if fi.IsDir() {
		lst, err := ioutil.ReadDir(absName)
		if err != nil {
			log.Fatalf("\t%s", err)
		}

		for _, item := range lst {
			execute(path.Join(absName, item.Name()), prefix)
		}
	}
}
