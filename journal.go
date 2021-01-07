package techlog

import (
	"context"
	"github.com/xelaj/go-dry"
	"io"
	"os"
	"regexp"
	"time"
)

var reHeaders = regexp.MustCompile(`(?m)(?P<Metadata>[0-9][0-9]:[0-9][0-9].[0-9]+-\d+,\w+,\d,)`)

type JournalFile struct {

	// yymmddhh.log
	// yy – две последние цифры года;
	// mm – номер месяца;
	// dd – номер дня;
	// hh – номер часа.
	File string
	Time time.Time

	endPos int64
}

func (j JournalFile) Watch(ctx context.Context) (update chan Event, err error) {

	update = make(chan Event)

	go func() {
		err = j.Parse(ctx, update)
	}()

	return

}

func (j JournalFile) Parse(ctx context.Context, update chan Event) error {

	ticker := time.NewTicker(30 * time.Second)

	free := make(chan struct{}, 1)
	free <- struct{}{}

	for {

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			<-free
			data, err := j.readFile()
			if err != nil {
				return err
			}

			if len(data) == 0 && time.Now().Sub(j.Time).Hours() >= 1 {
				ticker.Stop()
				return nil
			}

			j.readEvents(data, update)
			free <- struct{}{}
		}

	}

}

func (j *JournalFile) readEvents(data []byte, update chan Event) {

	headersIdx := reHeaders.FindAllIndex(data, -1)
	max := len(headersIdx)

	for i, idx := range headersIdx {

		startIdx := idx[0]
		endIdx := idx[1]

		header := data[startIdx:endIdx]

		endProps := len(data)
		if i+1 < max {
			endProps = headersIdx[i+1][0] - 1
		}

		props := data[endIdx:endProps]
		//idx := bytes.Index(data, match)

		e := newEvent(header, props)
		update <- e
	}
}

func (j *JournalFile) readFile() ([]byte, error) {

	openFile, err := os.Open(j.File)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 64)
	var data []byte

	for {
		n, err := openFile.ReadAt(buffer, j.endPos)
		j.endPos += int64(n)
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
		data = append(data, buffer[:n]...)
	}

	openFile.Close()
	return data, nil
}

func getFileTile(fileDate string) time.Time {

	year := dry.StringToInt(fileDate[0:2])
	month := dry.StringToInt(fileDate[2:4])
	day := dry.StringToInt(fileDate[4:6])
	hours := dry.StringToInt(fileDate[6:8])

	return time.Date(2000+year, time.Month(month), day, hours, 0, 0, 0, nil)

}
