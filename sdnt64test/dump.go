package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"
)

type interval struct {
	start time.Time
	end   time.Time
	delta int64
}

type intervals []interval

func makeReport(count int) intervals {
	return make(intervals, count)
}

func (r intervals) Len() int           { return len(r) }
func (r intervals) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r intervals) Less(i, j int) bool { return r[i].start.Before(r[j].start) }

func (r intervals) dump(start time.Time) {
	sort.Sort(r)

	w := bufio.NewWriter(os.Stdout)

	fmt.Fprint(w, "{\n  \"sends\": [")

	first := true
	for _, item := range r {
		n := item.start.Sub(start).Nanoseconds()
		if n < 0 {
			continue
		}

		if first {
			fmt.Fprintf(w, "\n    %d", n)
			first = false
		} else {
			fmt.Fprintf(w, ",\n    %d", n)
		}
	}

	fmt.Fprint(w, "\n  ],\n  \"receives\": [")

	recs := make([]int64, len(r))
	for i, item := range r {
		recs[i] = item.end.Sub(start).Nanoseconds()
	}
	sort.Sort(ints64(recs))

	first = true
	for _, n := range recs {
		if n < 0 {
			continue
		}

		if first {
			fmt.Fprintf(w, "\n    %d", n)
			first = false
		} else {
			fmt.Fprintf(w, ",\n    %d", n)
		}
	}

	fmt.Fprint(w, "\n  ],\n  \"pairs\": [")

	first = true
	for _, item := range r {
		n := item.start.Sub(start).Nanoseconds()
		if n < 0 {
			continue
		}

		if first {
			fmt.Fprintf(w, "\n    [%d", n)
			first = false
		} else {
			fmt.Fprintf(w, ",\n    [%d", n)
		}

		m := item.end.Sub(start).Nanoseconds()
		if m >= n {
			fmt.Fprintf(w, ", %d, %d", m, m-n)
		}

		fmt.Fprintf(w, "]")
	}
	fmt.Fprint(w, "\n  ]\n}\n")
	w.Flush()
}

type ints64 []int64

func (a ints64) Len() int           { return len(a) }
func (a ints64) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ints64) Less(i, j int) bool { return a[i] < a[j] }
