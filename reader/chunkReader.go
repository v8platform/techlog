package reader

import (
	"bufio"
	"io"
	"v8platform/techlog"
)

type chunkReader struct {
	rd           *bufio.Reader
	minChunkSize int
}

func NewChunkReader(r io.Reader, minChunkSize int) *chunkReader {
	return &chunkReader{
		rd:           bufio.NewReader(r),
		minChunkSize: minChunkSize,
	}
}

// Выполняет чтение из редера до ближайшего окончания лога
func (r *chunkReader) Read() ([]byte, int, error) {

	buf := make([]byte, r.minChunkSize)

	n, err := r.rd.Read(buf)
	if err != nil {
		return nil, n, err
	}

	if n < len(buf) {
		buf = buf[:n]
	}

	err = r.readForNextLog(&buf)

	if err != nil && err != io.EOF {
		return nil, 0, err
	}

	return buf, len(buf), err

}

func (r *chunkReader) readForNextLog(buf *[]byte) error {

	var findBuffer int
	var n int

	findBuffer = 512
	maxBufSize := r.minChunkSize
	size := 0
	for {
		n++
		peekSize := n * findBuffer
		fbuf, err := r.rd.Peek(peekSize)
		if err != nil {
			size = len(fbuf)
			break
		}

		idx := techlog.reHeaders.FindIndex(fbuf)

		if idx == nil {

			if peekSize >= maxBufSize {
				add := make([]byte, peekSize)
				_, err := r.rd.Read(add)
				if err != nil {
					return err
				}
				*buf = append(*buf, add...)
				n = 0
			}

			continue
		}

		size = idx[0]
		break
	}

	if size == 0 {
		size = maxBufSize
	}

	add := make([]byte, size)
	_, err := r.rd.Read(add)
	if err != nil {
		return err
	}
	*buf = append(*buf, add...)

	return nil
}
