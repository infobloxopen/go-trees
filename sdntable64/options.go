package sdntable64

type Option func(o *options)

func WithPath(path string) Option {
	return func(o *options) {
		o.path = &path
	}
}

func WithLimit(limit int) Option {
	return func(o *options) {
		o.limit = limit
	}
}

func WithNormalizeLoggers(before, sort, after func(size, length int)) Option {
	return func(o *options) {
		o.log.normalize.before = before
		o.log.normalize.sort = sort
		o.log.normalize.after = after
	}
}

func WithFlushLoggers(before, after func(size, from, to int)) Option {
	return func(o *options) {
		o.log.flush.before = before
		o.log.flush.after = after
	}
}

func WithReadLogger(log func(size, from, to, reqs, queue int)) Option {
	return func(o *options) {
		o.log.read = log
	}
}

type options struct {
	path  *string
	limit int
	log   loggers
}
