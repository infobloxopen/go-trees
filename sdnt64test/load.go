package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func load(path string, header func(n int) error, pair func(k string, v uint64) error) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	i := 0
	n := -1
	for s.Scan() {
		t := s.Text()
		if n < 0 {
			v, err := strconv.ParseInt(t, 0, 32)
			if err != nil {
				return fmt.Errorf("can't treat %q as number of domains: %s", t, err)
			}

			if v < 0 {
				return fmt.Errorf("negative number of lines: %d", v)
			}

			n = int(v)
			if err := header(n); err != nil {
				return err
			}
		} else {
			i++
			if i > n {
				return fmt.Errorf("got more lines %d than expected %d", i, n)
			}

			flds := strings.SplitN(t, " ", 2)
			if len(flds) != 2 {
				return fmt.Errorf("expected two fields (categories and domain) but got %d: %q", len(flds), t)
			}

			v, err := strconv.ParseUint(flds[0], 16, 64)
			if err != nil {
				return fmt.Errorf("can't treat %q as categories: %s", flds[0], err)
			}

			k, err := strconv.Unquote(flds[1])
			if err != nil {
				return fmt.Errorf("can't treat %q as quoted string: %s", flds[1], err)
			}

			if err := pair(k, v); err != nil {
				return err
			}
		}
	}

	if err := s.Err(); err != nil {
		return err
	}

	if i != n {
		return fmt.Errorf("expected %d domains but got %d", n, i)
	}

	return nil
}
