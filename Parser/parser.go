package Parser

import (
	"bytes"
	"github.com/xelaj/go-dry"
	"regexp"
	"strings"
	"time"
	"v8platform/techlog"
)

type Parser interface {
	Parse(data []byte, t time.Time, offset int64) []techlog.Event
	HeaderIndex(data []byte) []int
}

func NewDataParser(reHeaders, reGaps *regexp.Regexp) Parser {
	return parser{
		reHeaders: reHeaders,
		reGaps:    reGaps,
	}
}

type parser struct {
	reHeaders *regexp.Regexp
	reGaps    *regexp.Regexp
}

func (p parser) HeaderIndex(data []byte) []int {
	return p.reHeaders.FindIndex(data)
}

func (p parser) Parse(data []byte, t time.Time, offset int64) []techlog.Event {

	headersIdx := p.reHeaders.FindAllIndex(data, -1)
	max := len(headersIdx)
	var events []techlog.Event
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

		e := techlog.Event{
			Time:   t,
			Offset: off,
			Size:   len(header) + len(props),
		}
		e.TimeOffset, e.Duration, e.Type, e.StackLevel = p.parseEventHeader(header)

		if len(props) > 0 {
			e.Properties = p.parseEventProps(props)
		}

		events = append(events, e)

	}

	return events
}

func (p parser) parseEventHeader(header []byte) (time.Duration, time.Duration, techlog.EventType, int) {

	header = bytes.TrimRight(header, ",")
	data := bytes.Split(header, []byte(","))

	StackLevel := dry.StringToInt(string(data[2]))
	Type := techlog.EventType(data[1])

	logTime := bytes.Split(data[0], []byte("-"))

	min := int64(dry.StringToInt(string(logTime[0][0:2]))) * int64(time.Minute)
	sec := int64(dry.StringToInt(string(logTime[0][3:5]))) * int64(time.Second)
	nsec := int64(dry.StringToInt(string(logTime[0][6:]))) * int64(time.Microsecond)

	TimeOffset := time.Duration(min + sec + nsec)
	Duration := time.Duration(dry.StringToInt(string(logTime[1])))

	return TimeOffset, Duration, Type, StackLevel

}

func (p parser) parseEventProps(rawData []byte) map[string]string {

	rows := bytes.Split(p.replaceGaps(rawData), []byte(","))

	props := make(map[string]string, len(rows))

	for _, v := range rows {
		if len(v) == 0 {
			continue
		}

		keyValue := bytes.SplitN(v, []byte("="), 2)
		//dry.PanicIf(len(keyValue) != 2, "error parse props ", string(rawData))
		key := strings.Replace(string(keyValue[0]), ":", "_", 1)

		value := string(p.restoreGaps(keyValue[1]))
		props[key] = value

	}

	return props
}

// заменяет , на пробел в строках вида ' , ,  '
func (p parser) replaceGaps(data []byte) []byte {
	data = bytes.TrimRight(data, ",")
	data = bytes.TrimSpace(data)

	gapsStrings := p.reGaps.FindAll(data, -1)
	for _, gapString := range gapsStrings {
		data = bytes.Replace(data, gapString, bytes.ReplaceAll(gapString, []byte(","), []byte{0x7F}), -1)
	}

	return data
}

// заменяет , на пробел в строках вида ' , ,  '
func (p parser) restoreGaps(data []byte) []byte {

	data = bytes.ReplaceAll(data, []byte{0x7F}, []byte(","))

	return data
}
