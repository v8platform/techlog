package techlog

import (
	"bytes"
	"sync"
)

var _bufferPool = sync.Pool{New: func() interface{} {
	return &bytes.Buffer{}
}}

func GetBuffer() *bytes.Buffer {
	buf := _bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	_bufferPool.Put(buf)
}
