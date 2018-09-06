package main

import "log"

type dMap struct {
	m map[string]uint64
}

func newMap(in []pair) *dMap {
	out := &dMap{
		m: make(map[string]uint64, len(in)),
	}

	for i, p := range in {
		d, err := domainUppercase(p.k)
		if err != nil {
			log.Fatalf("can't make case insensitive domain from %q: %s", p.k, err)
		}

		out.m[d] = p.v

		if (i+1)%1000000 == 0 {
			log.Printf("inserted %d domains", i+1)
		}
	}

	if len(in)%1000000 != 0 {
		log.Printf("inserted %d domains", len(in))
	}

	return out
}

func (m *dMap) Map(k string) uint64 {
	d, err := domainUppercase(k)
	if err != nil {
		log.Fatalf("can't make case insensitive domain from %q: %s", k, err)
	}

	for {
		if v, ok := m.m[d]; ok {
			return v
		}

		if len(d) <= 0 {
			break
		}

		d, err = dropLabel(d)
		if err != nil {
			log.Fatalf("can't extract next level zone from %q: %s", k, err)
		}
	}

	return 0
}
