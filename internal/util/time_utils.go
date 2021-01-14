package util

import (
	"time"
)

func CurrentTimestampMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}