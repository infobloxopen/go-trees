package main

import (
	"fmt"
	"log"

	"github.com/infobloxopen/go-trees/udomain"
)

func generate(d [127][]int64, b, m, a int) ([]string, error) {
	out := []string{}

	for i, d := range d {
		if len(d) <= 0 {
			continue
		}

		log.Printf("generating domains for %d (%d)", i+1, len(d)/i+1)
		s, err := genMissing(d, i+1, b, m, a)
		if err != nil {
			return nil, err
		}

		out = append(out, s...)
	}

	return out, nil
}

func genMissing(d []int64, n, b, m, a int) ([]string, error) {
	out := []string{}

	if len(d) > 0 {
		s, err := genLeading(d[:n], b)
		if err != nil {
			return nil, err
		}
		out = append(out, s...)

		s, err = genInternal(d, n, m)
		if err != nil {
			return nil, err
		}
		out = append(out, s...)

		s, err = genFollowing(d[len(d)-n:], a)
		if err != nil {
			return nil, err
		}
		out = append(out, s...)
	}

	return out, nil
}

func genLeading(min []int64, n int) ([]string, error) {
	dn, err := domain.MakeNameFromSlice(min)
	if err != nil {
		return nil, fmt.Errorf("convert min %d: %s", len(min), err)
	}

	out := make([]string, 0, n)
	for j := 0; j < n; j++ {
		p, err := prevDomain(min)
		if err != nil {
			return nil, fmt.Errorf("generate leading %q (%d): %s", dn, len(min), err)
		}

		tmp, err := domain.MakeNameFromSlice(p)
		if err != nil {
			return nil, fmt.Errorf("convert leading %q (%d): %s", dn, len(min), err)
		}

		min = p
		dn = tmp

		out = append(out, dn.String())
	}

	return out, nil
}

func genFollowing(max []int64, n int) ([]string, error) {
	dn, err := domain.MakeNameFromSlice(max)
	if err != nil {
		return nil, fmt.Errorf("convert max %d: %s", len(max), err)
	}

	out := make([]string, 0, n)
	for j := 0; j < n; j++ {
		next, err := nextDomain(max)
		if err != nil {
			return nil, fmt.Errorf("generate following %q (%d): %s", dn, len(max), err)
		}

		tmp, err := domain.MakeNameFromSlice(next)
		if err != nil {
			return nil, fmt.Errorf("convert following %q (%d): %s", dn, len(max), err)
		}

		max = next
		dn = tmp

		out = append(out, dn.String())
	}

	return out, nil
}

func genInternal(d []int64, n, m int) ([]string, error) {
	out := []string{}

	i := 0
	for j := 0; j+n <= len(d); j += n * m {
		if j-n < 0 {
			continue
		}

		min := d[j-n : j]
		mnDn, err := domain.MakeNameFromSlice(min)
		if err != nil {
			return nil, fmt.Errorf("convert left %d at %d (%d:%d): %s", n, i+1, j-n, j, err)
		}

		max := d[j : j+n]
		mxDn, err := domain.MakeNameFromSlice(max)
		if err != nil {
			return nil, fmt.Errorf("convert right %d at %d (%d:%d): %s", n, i+1, j, j+n, err)
		}

		p, err := prevDomain(max)
		if err != nil {
			return nil, fmt.Errorf("generate inside %q - %q (%d) at %d (%d:%d): %s",
				mnDn, mxDn, n, i+1, j-n, j+n, err)
		}

		dn, err := domain.MakeNameFromSlice(p)
		if err != nil {
			return nil, fmt.Errorf("convert inside %q - %q (%d) at %d (%d:%d): %s",
				mnDn, mxDn, n, i+1, j-n, j+n, err)
		}

		if compare(min, p) < 0 {
			out = append(out, dn.String())

			i++
			if i%100000 == 0 {
				log.Printf("generated %d domains", i)
			}
		} else {
			log.Printf("no space inside %q - %q: %q (%d) at %d (%d:%d)", mnDn, mxDn, dn, n, i+1, j-n, j+n)
		}
	}

	if i%100000 != 0 {
		log.Printf("generated %d domains", i)
	}

	return out, nil
}
