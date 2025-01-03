package utils

import (
	"time"
)

func GetCurrentTimeStampForIndia() time.Time {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	return time.Now().In(loc)
}

func GetCurrentEpochTimeStamp() int64 {
	return time.Now().In(time.UTC).UnixMilli()
}
