package reader

import (
	"v8platform/techlog/Parser"
)

type Options struct {
	minChunkSize    int
	parallelReaders int
	offset          int64
	dataParser      Parser.Parser
}

func WithChunkSize(size int) func(opt *Options) {
	return func(opt *Options) {
		opt.minChunkSize = size
	}
}
func WithParallels(poolSize int) func(opt *Options) {
	return func(opt *Options) {
		opt.parallelReaders = poolSize
	}
}
func WithOffset(offset int64) func(opt *Options) {
	return func(opt *Options) {
		opt.offset = offset
	}
}
func WithParser(dataParser Parser.Parser) func(opt *Options) {
	return func(opt *Options) {
		opt.dataParser = dataParser
	}
}
