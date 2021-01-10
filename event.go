package techlog

import (
	"bytes"
	"github.com/xelaj/go-dry"
	"regexp"
	"strings"
	"time"
)

type EventType string

const (
	CallType EventType = "CALL"
)

var regexGaps = regexp.MustCompile(`(?m)('[\S\s]*?')|("[\S\s]*?")`)
var reHeaders = regexp.MustCompile(`(?m)(?P<Metadata>[0-9][0-9]:[0-9][0-9].[0-9]+-\d+,\w+,\d,)`)

type Event struct {

	// Дата и время события до часа из даты файла
	Time time.Time

	// Момент события от даты файла
	TimeOffset time.Duration

	// Отспут до блока события в файле
	Offset int64

	// Размер блока события в файле
	Size int

	// Длительность выполнения в милисикундах
	Duration time.Duration

	// Тип события
	Type EventType

	// Уровень события в стеке выполнения
	StackLevel int

	// Набор свойств события
	Properties map[string]string
}

// Дата начала события (расчетная, т.к. в журнале есть ктолько окончания)
func (e Event) StartAt(t ...time.Time) time.Time {
	return e.EndAt(t...).Add(-e.Duration)
}

// Дата окончания события
func (e Event) EndAt(t ...time.Time) time.Time {

	tt := e.Time

	if len(t) > 0 {
		tt = t[0]
	}

	return tt.Add(e.TimeOffset)
}

func parseChunkData(data []byte, t time.Time, offset int64) []Event {

	//str := string(data)
	headersIdx := reHeaders.FindAllIndex(data, -1)
	max := len(headersIdx)
	var events []Event
	for i, idx := range headersIdx {

		startIdx := idx[0]
		endIdx := idx[1]

		header := data[startIdx:endIdx]

		endProps := len(data)
		if i+1 < max {
			endProps = headersIdx[i+1][0] - 1
		}

		props := data[endIdx:endProps]

		off := int64(startIdx) + offset

		e := Event{
			Time:   t,
			Offset: off,
			Size:   len(header) + len(props),
		}
		e.TimeOffset, e.Duration, e.Type, e.StackLevel = parseEventHeader(header)

		if len(props) > 0 {
			e.Properties = parseEventProps(props)
		}

		events = append(events, e)

	}

	return events

}

func parseEventHeader(header []byte) (time.Duration, time.Duration, EventType, int) {

	header = bytes.TrimRight(header, ",")
	data := bytes.Split(header, []byte(","))

	StackLevel := dry.StringToInt(string(data[2]))
	Type := EventType(data[1])

	logTime := bytes.Split(data[0], []byte("-"))

	min := int64(dry.StringToInt(string(logTime[0][0:2]))) * int64(time.Minute)
	sec := int64(dry.StringToInt(string(logTime[0][3:5]))) * int64(time.Second)
	nsec := int64(dry.StringToInt(string(logTime[0][6:]))) * int64(time.Microsecond)

	TimeOffset := time.Duration(min + sec + nsec)
	Duration := time.Duration(dry.StringToInt(string(logTime[1])))

	return TimeOffset, Duration, Type, StackLevel

}

func parseEventProps(rawData []byte) map[string]string {

	p := bytes.Split(replaceGaps(rawData), []byte(","))

	props := make(map[string]string, len(p))

	for _, v := range p {
		if len(v) == 0 {
			continue
		}

		keyValue := bytes.SplitN(v, []byte("="), 2)
		//dry.PanicIf(len(keyValue) != 2, "error parse props ", string(rawData))
		key := strings.Replace(string(keyValue[0]), ":", "_", 1)

		value := string(restoreGaps(keyValue[1]))
		props[key] = value

	}

	return props

}

// заменяет , на пробел в строках вида ' , ,  '
func replaceGaps(data []byte) []byte {
	data = bytes.TrimRight(data, ",")
	data = bytes.TrimSpace(data)

	gapsStrings := regexGaps.FindAll(data, -1)
	for _, gapString := range gapsStrings {
		data = bytes.Replace(data, gapString, bytes.ReplaceAll(gapString, []byte(","), []byte{0x7F}), -1)
	}

	return data
}

// заменяет , на пробел в строках вида ' , ,  '
func restoreGaps(data []byte) []byte {

	data = bytes.ReplaceAll(data, []byte{0x7F}, []byte(","))

	return data
}
