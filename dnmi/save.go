package main

import (
	"bufio"
	"fmt"
	"os"
)

func save(path string, s []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	if _, err := fmt.Fprintf(w, "%d\n", len(s)); err != nil {
		return fmt.Errorf("failed to write header to %q: %s", path, err)
	}

	for i, s := range s {
		if _, err := fmt.Fprintf(w, "%q\n", s); err != nil {
			return fmt.Errorf("failed to writed domain %q (%d) to %q: %s", s, i+1, path, err)
		}
	}

	return w.Flush()
}
