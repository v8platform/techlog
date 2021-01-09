package techlog

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type EventType string

const (
	CallType EventType = "CALL"
)

type Event struct {

	// Дата начала события (расчетная, т.к. в журнале есть ктолько окончания)
	StartAt time.Time
	// Дата окончания события
	EndAt time.Time
	// Длительность выполнения в милисикундах
	Duration time.Duration
	// Тип события
	Type EventType
	// Уровень события в стеке выполнения
	StackLevel string

	// Набор свойств события
	Props map[string]string
}

func parseChunkData(data []byte, t time.Time) []Event {

	var reHeaders = regexp.MustCompile(`(?m)(?P<Metadata>[0-9][0-9]:[0-9][0-9].[0-9]+-\d+,\w+,\d,)`)
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
		//idx := bytes.Index(data, match)

		e := newEvent(t, header, props)
		events = append(events, e)

		//pp.Println("event", e)
	}

	return events

}

func newEvent(t time.Time, header []byte, props ...[]byte) Event {

	e := Event{}

	header = bytes.TrimRight(header, ",")
	data := bytes.Split(header, []byte(","))

	e.StackLevel = string(data[2])
	e.Type = EventType(data[1])

	logTime := bytes.Split(data[0], []byte("-"))
	min, _ := strconv.ParseInt(string(logTime[0][0:2]), 10, 64)
	sec, _ := strconv.ParseInt(string(logTime[0][3:5]), 10, 64)
	nsec, _ := strconv.ParseInt(string(logTime[0][6:]), 10, 64)

	minutes := int64(time.Minute) * min
	seconds := int64(time.Second) * sec
	nseconds := int64(time.Microsecond) * nsec
	e.EndAt = t.Add(time.Duration(minutes + seconds + nseconds))

	d, _ := strconv.ParseInt(string(logTime[1]), 10, 64)

	e.Duration = time.Duration(d)
	e.StartAt = e.EndAt.Add(-e.Duration * time.Millisecond)
	if len(props) > 0 {
		e.Props = parseEventProps(props[0])
	}

	return e
}

func parseEventProps(rawData []byte) map[string]string {

	props := make(map[string]string)

	p := bytes.Split(replaceGaps(rawData), []byte(","))

	for _, v := range p {
		if len(v) == 0 {
			continue
		}

		keyValue := bytes.SplitN(v, []byte("="), 2)
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
	regexGaps := regexp.MustCompile("('{1}[\\S\\s]*?'{1})|((\"{1}[\\S\\s]*?\"{1}))")
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
