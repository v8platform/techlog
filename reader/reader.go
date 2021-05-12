package reader

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"time"
	"v8platform/techlog"
	"v8platform/techlog/Parser"
	"v8platform/techlog/utils"
)

type Reader interface {
	Read() ([]techlog.Event, int64, error)
}

type StreamReader interface {
	Read(in chan techlog.Event) (int64, error)
}

var defOptions = Options{
	minChunkSize:    4 * 1024,
	parallelReaders: 0,
	dataParser: Parser.NewDataParser(
		regexp.MustCompile(`(?m)(?P<Metadata>[0-9][0-9]:[0-9][0-9].[0-9]+-\d+,\w+,\d,)`),
		regexp.MustCompile(`(?m)('[\S\s]*?')|("[\S\s]*?")`),
	),
}

func NewReader(r io.Reader, time time.Time, opts ...func(options *Options)) Reader {

	options := defOptions
	for _, opt := range opts {
		opt(&options)
	}

	return &reader{
		rd:        bufio.NewReader(r),
		eventDate: time,
	}

}

func NewFileReader(file *os.File, offset int64, opts ...func(options *Options)) (Reader, error) {

	options := defOptions
	options.offset = offset
	for _, opt := range opts {
		opt(&options)
	}

	filestats, _ := file.Stat()

	t := utils.GetFileDatetime(filestats.Name())

	if options.offset > 0 {
		_, err := file.Seek(offset, io.SeekStart)
		if err != nil {
			return nil, err
		}
	}

	return &reader{
		rd:        bufio.NewReader(file),
		eventDate: t,
	}, nil

}

type reader struct {
	Options
	rd        *bufio.Reader
	eventDate time.Time
	chunk     []byte
	chunkSize int
}

func (r *reader) Read() ([]techlog.Event, int64, error) {

	var readBytes int64
	var events []techlog.Event

	for {

		data, n, err := r.readChunk()

		if n > 0 {
			newEvents := r.dataParser.Parse(data, r.eventDate, r.offset)
			events = append(events, newEvents...)
		}

		r.offset += int64(n)
		readBytes += int64(n)

		if err == io.EOF {
			return events, readBytes, err
		}

		if err != nil && err != io.EOF {
			return events, readBytes, err
		}

	}
}

// Выполняет чтение из редера до ближайшего окончания лога
func (r *reader) readChunk() ([]byte, int, error) {

	buf := make([]byte, r.minChunkSize)

	n, err := r.rd.Read(buf)
	if err != nil {
		return nil, n, err
	}

	if n < len(buf) {
		buf = buf[:n]
	}

	err = r.appendEOFChunk(&buf)

	if err != nil && err != io.EOF {
		return nil, 0, err
	}

	return buf, len(buf), err

}

func (r *reader) appendEOFChunk(buf *[]byte) error {

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

		idx := r.dataParser.HeaderIndex(fbuf)

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
