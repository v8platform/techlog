package techlog

import (
	"bufio"
	"context"
	"github.com/radovskyb/watcher"
	"github.com/xelaj/go-dry"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Options func()
type Events chan Event

var DefaultChunkSize = 4 * 1024

func MaxEvents(max int) Options {
	return func() {

	}
}

func FilterType(filters ...string) Options {
	return func() {

	}
}

type Watcher struct {
	dir string

	MaxEvents   int
	Filters     []string
	Ctx         context.Context
	e           Events
	journals    map[string]int64
	fileWatcher *watcher.Watcher
}

func (w *Watcher) Stop() {
	// Завершает конекст выполнения чтения
	w.Ctx.Done()
	w.fileWatcher.Close()
}

func (w *Watcher) Start() Events {

	w.e = make(Events, w.MaxEvents)

	w.process()

	return w.e
}

func (w *Watcher) process() {

	go w.processFileWatcher()

}

func (w *Watcher) processFileWatcher() {

	w.journals = make(map[string]int64)

	w.fileWatcher = watcher.New()
	if err := w.fileWatcher.Add(w.dir); err != nil {
		log.Fatalln(err)
	}
	w.fileWatcher.SetMaxEvents(1)
	w.fileWatcher.FilterOps(watcher.Create, watcher.Write, watcher.Remove)

	go func() {
		_ = w.fileWatcher.Start(time.Microsecond * 300)
	}()

	for {
		select {
		case <-w.Ctx.Done():
			break
		case e := <-w.fileWatcher.Event:
			switch e.Op {

			case watcher.Write:
				w.emmitWriteEvent(e)
			case watcher.Create:

				w.journals[e.Name()] = 0

			case watcher.Remove:
				delete(w.journals, e.Name())
			}

		}
	}
}

func (w *Watcher) emmitWriteEvent(event watcher.Event) {

	if event.IsDir() || strings.HasSuffix(event.Name(), "~") {
		return
	}

	offset := w.journals[event.Name()]

	newOffset, _ := readLogFile(event.Path, offset, w.e)

	w.journals[event.Name()] = newOffset

}

func readLogFile(path string, offset int64, inEvents Events) (n int64, err error) {

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)

	if err != nil {
		return 0, err
	}
	filestats, _ := file.Stat()

	t := getFileDatetime(filestats.Name())

	if offset > 0 {
		_, err = file.Seek(offset, io.SeekStart)
		if err != nil {
			return
		}
	}
	n, err = readTechlogData(file, offset, t, inEvents)

	if err != nil {
		return n, err
	}

	return offset + n, nil
}

func readTechlogData(reader io.Reader, offset int64, t time.Time, in Events) (int64, error) {

	var readBytes int64

	wg := &sync.WaitGroup{}
	//eLock := &sync.Mutex{}
	limitReader := make(chan struct{}, 20)

	cr := NewChunkReader(reader, DefaultChunkSize)

	for {
		limitReader <- struct{}{}
		data, n, err := cr.Read()

		switch err {
		case nil, io.EOF:
			//
		default:
			log.Printf("error reading data <%s>", err)
			<-limitReader
			break
		}
		if n == 0 {
			<-limitReader
			break
		}

		wg.Add(1)
		go func(d []byte, off int64) {
			events := parseChunkData(d, t, off)
			//eLock.Lock()
			for _, event := range events {
				in <- event
			}
			//eLock.Unlock()
			wg.Done()
			<-limitReader
		}(data, offset)

		offset += int64(n)
		readBytes += int64(n)

	}

	wg.Wait()

	return readBytes, nil

}

func getFileDatetime(name string) time.Time {

	year := "20" + name[0:2]
	month := name[2:4]
	day := name[4:6]
	hours := name[6:8]

	return time.Date(dry.StringToInt(year), time.Month(dry.StringToInt(month)),
		dry.StringToInt(day), dry.StringToInt(hours), 0, 0, 0, time.Local)

}

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
func (cr *chunkReader) Read() ([]byte, int, error) {

	buf := make([]byte, cr.minChunkSize)

	n, err := cr.rd.Read(buf)
	if err != nil {
		return nil, n, err
	}

	if n < len(buf) {
		buf = buf[:n]
	}

	err = cr.readForNextLog(&buf)

	if err != nil && err != io.EOF {
		return nil, 0, err
	}

	return buf, len(buf), err

}

func (cr *chunkReader) readForNextLog(buf *[]byte) error {

	var findBuffer int
	var n int

	findBuffer = 512
	maxBufSize := 4096
	size := 0
	for {
		n++
		peekSize := n * findBuffer
		fbuf, err := cr.rd.Peek(peekSize)
		if err != nil {
			return err
		}

		idx := reHeaders.FindIndex(fbuf)

		if idx == nil {

			if peekSize >= maxBufSize {
				add := make([]byte, peekSize)
				_, err := cr.rd.Read(add)
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

	add := make([]byte, size)
	_, err := cr.rd.Read(add)
	if err != nil {
		return err
	}
	*buf = append(*buf, add...)

	return nil
}

// TODO Пока сделал заготовку для функции мониторинга файла
func watch(dir string, opts ...Options) (Events, *Watcher) {

	w := &Watcher{
		dir:       dir,
		MaxEvents: 10,
		Ctx:       context.Background(),
	}

	return w.Start(), w

}

func StreamRead(file string, maxEvents int64) (Events, error) {

	events, _, err := StreamReadAt(file, maxEvents, 0)
	return events, err

}

func StreamReadAt(file string, maxEvents int64, offset int64) (Events, *int64, error) {

	fullPath, err := filepath.Abs(file)
	if err != nil {
		return nil, &offset, err
	}

	off := &offset

	stream := make(Events, maxEvents)

	go func() {
		*off, err = readLogFile(fullPath, offset, stream)
		close(stream)
	}()

	return stream, off, nil

}

func ReadAt(file string, offset int64) ([]Event, *int64, error) {

	stream, off, err := StreamReadAt(file, 50, offset)

	if err != nil {
		return nil, &offset, err
	}

	var events []Event

	for event := range stream {
		events = append(events, event)
	}

	return events, off, nil

}

func Read(file string) ([]Event, error) {

	events, _, err := ReadAt(file, 0)
	return events, err
}
