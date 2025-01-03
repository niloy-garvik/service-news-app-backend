package utils

import (
	"fmt"
	"time"
)

func GetCurrentTimeStampForIndia() time.Time {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	return time.Now().In(loc)
}

func GetCurrentEpochTimeStamp() int64 {
	return time.Now().In(time.UTC).UnixMilli()
}

// ConvertStringToTimestamp converts a string in RFC3339 format to time.Time
func ConvertStringToTimestamp(dateString string) (time.Time, error) {
	// Parse the date string using the RFC3339 layout
	layout := time.RFC3339
	timestamp, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date string: %v", err)
	}
	return timestamp, nil
}
