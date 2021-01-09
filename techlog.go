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
	"regexp"
	"strings"
	"sync"
	"time"
)

type Options func()
type Events chan Event

var ChunkReaderSize = 4 * 1024

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

	newOffset := readLogFile(event.Path, offset, w.e)

	w.journals[event.Name()] = newOffset

}

func readLogFile(path string, offset int64, inEvents Events) int64 {

	cr := chunkReader{ChunkReaderSize}

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)

	if err != nil {
		log.Panicf("emmit write event: open file %s", err)
	}
	filestats, _ := file.Stat()

	t := getFileDatetime(filestats.Name())

	wg := &sync.WaitGroup{}
	eLock := &sync.Mutex{}
	for {
		chunk := cr.Read(file, offset)
		if chunk.size == 0 {
			break
		}
		wg.Add(1)
		go func(eData []byte) {
			events := parseChunkData(chunk.data, t)
			eLock.Lock()
			for _, event := range events {
				inEvents <- event
			}
			eLock.Unlock()
			wg.Done()
		}(chunk.data)
		offset = chunk.end
	}

	wg.Wait()

	return offset
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
	chunkSize int
}

type chunk struct {
	data  []byte
	start int64
	size  int
	end   int64
}

var reHeaders = regexp.MustCompile(`(?mi)([0-9][0-9]:[0-9][0-9]\.[0-9]+-\d+)`)

func (cr chunkReader) Read(file *os.File, offset int64) chunk {

	//file, _ := os.Open(j.File)
	//defer file.Close()
	buf := make([]byte, cr.chunkSize)

	var err error
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return chunk{}
	}
	n, _ := file.Read(buf)
	if n == 0 {
		return chunk{}
	}

	c := chunk{
		start: offset,
		size:  n,
	}

	if c.size < cr.chunkSize {
		c.data = buf[:c.size]
		c.end = c.start + int64(c.size)
		return c
	}

	reader := bufio.NewReader(file)

	for {
		txt, err := reader.ReadBytes('\n')

		if ok := reHeaders.Match(txt); ok {
			break
		}

		buf = append(buf, txt...)
		c.size += len(txt)

		if err == io.EOF {
			break
		}

	}

	c.end = c.start + int64(c.size)
	c.data = buf

	return c

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

func StreamReadAt(file string, maxEvents int64, offset int64) (Events, int64, error) {

	fullPath, err := filepath.Abs(file)
	if err != nil {
		return nil, offset, err
	}

	stream := make(Events, maxEvents)

	go func() {
		offset = readLogFile(fullPath, offset, stream)
		close(stream)
	}()

	return stream, offset, nil

}

func ReadAt(file string, offset int64) ([]Event, int64, error) {

	stream, offset, err := StreamReadAt(file, 50, offset)

	if err != nil {
		return nil, offset, err
	}

	var events []Event

	for event := range stream {
		events = append(events, event)
	}

	return events, offset, nil

}

func Read(file string) ([]Event, error) {

	events, _, err := ReadAt(file, 0)
	return events, err
}
