package sdntable64

type loggers struct {
	normalize normalizeLoggers
	flush     flushLoggers
	read      func(size, from, to, reqs, queue int)
}

type normalizeLoggers struct {
	before func(size, length int)
	sort   func(size, length int)
	after  func(size, length int)
}

type flushLoggers struct {
	before func(size, from, to int)
	after  func(size, from, to int)
}
