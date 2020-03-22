package utils

import (
	"time"
)

func GetTimeNowFormatYYYYMMDDHHIIMM() string {
	t := time.Now()
	s := t.Format("20060102150405")
	return s
}
