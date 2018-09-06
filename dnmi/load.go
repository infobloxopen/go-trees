package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func load(path string) ([]string, error) {
	var out []string

	f, err := os.Open(path)
	if err != nil {
		return out, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	i := 0
	for s.Scan() {
		t := s.Text()
		if len(out) <= 0 {
			v, err := strconv.ParseInt(t, 0, 32)
			if err != nil {
				return out, fmt.Errorf("can't treat %q as number of domains: %s", t, err)
			}

			if v < 0 {
				return out, fmt.Errorf("negative number of lines: %d", v)
			}

			out = make([]string, int(v))
		} else {
			if i >= len(out) {
				return out, fmt.Errorf("got more lines %d than expected %d", i+1, len(out))
			}

			flds := strings.SplitN(t, " ", 2)
			if len(flds) != 2 {
				return out, fmt.Errorf("expected two fields (categories and domain) but got %d: %q", len(flds), t)
			}

			k, err := strconv.Unquote(flds[1])
			if err != nil {
				return out, fmt.Errorf("can't treat %q as quoted string: %s", flds[1], err)
			}

			out[i] = k
			i++
		}
	}

	if err := s.Err(); err != nil {
		return out, err
	}

	if i != len(out) {
		return out, fmt.Errorf("expected %d domains but got %d", len(out), i)
	}

	return out, nil
}
