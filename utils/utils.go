package utils

import (
	"github.com/xelaj/go-dry"
	"time"
)

func GetFileDatetime(name string) time.Time {

	year := "20" + name[0:2]
	month := name[2:4]
	day := name[4:6]
	hours := name[6:8]

	return time.Date(dry.StringToInt(year), time.Month(dry.StringToInt(month)),
		dry.StringToInt(day), dry.StringToInt(hours), 0, 0, 0, time.Local)

}
