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
		out.m[p.k] = p.v

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
	var err error
	for {
		if v, ok := m.m[k]; ok {
			return v
		}

		if len(k) <= 0 {
			break
		}

		k, err = dropLabel(k)
		if err != nil {
			log.Fatalf("can't extract next level zone from %q: %s", k, err)
		}
	}

	return 0
}
