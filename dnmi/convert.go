package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/infobloxopen/go-trees/udomain"
)

func convert(ss []string) ([127][]int64, error) {
	var out [127][]int64
	n := 0
	for i, s := range ss {
		d, err := domain.MakeNameFromString(s)
		if err != nil {
			return out, fmt.Errorf("invalid domain %q at %d: %s", s, i+1, err)
		}

		c := d.GetComparable()
		if len(c) > 0 {
			j := len(c) - 1

			a := out[j]
			if len(a) <= 0 {
				a = []int64{}
			}

			out[j] = append(a, c...)

			n++
			if n%1000000 == 0 {
				log.Printf("converted %d domains", n)
			}
		}
	}

	if n%1000000 != 0 {
		log.Printf("converted %d domains", n)
	}

	for i, a := range out {
		if len(a) > 0 {
			s := makeSortable(a, i+1)
			sort.Sort(s)
			out[i] = s.getSorted()
		}
	}

	return out, nil
}
