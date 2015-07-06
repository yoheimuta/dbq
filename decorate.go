package main

import (
	"fmt"
	"regexp"
	"time"
)

var tableRule = regexp.MustCompile("@")

func Decorate(statement string, hour float64, start string, end string, hadd float64, buffer float64) (decorated string, err error) {
	if start != "" {
		decorated, err = withDateTime(statement, start, end, hadd, buffer)
	} else {
		if hour <= 0 {
			hour = 0.5
		}
		decorated, err = withHour(statement, hour)
	}
	return decorated, err
}

func GetRawStatement(statement string) (raw string) {
	// input:  tableName@
	// output: tableName
	raw = tableRule.ReplaceAllString(statement, "")
	return raw
}

func withDateTime(statement string, start string, end string, hadd float64, buffer float64) (decorated string, err error) {
	startTime, err := time.Parse("2006-01-02 15:04:05", start)
	if err != nil {
		return "", err
	}
	startTime = addHour(startTime, hadd)
	if 0 < buffer {
		startTime = addHour(startTime, buffer*-1)
	}
	startMSec := startTime.Unix() * 1000

	var replaced string
	if end == "" {
		replaced = fmt.Sprintf("@%d-", startMSec)
	} else {
		endTime, err := time.Parse("2006-01-02 15:04:05", end)
		if err != nil {
			return "", err
		}
		endTime = addHour(endTime, hadd)
		if 0 < buffer {
			endTime = addHour(endTime, buffer)
		}
		endMSec := endTime.Unix() * 1000
		replaced = fmt.Sprintf("@%d-%d", startMSec, endMSec)
	}

	// input:  tableName@
	// output: tableName@1435997334864-1436001613000
	decorated = tableRule.ReplaceAllString(statement, replaced)

	if statement == decorated {
		return "", fmt.Errorf("Failed to decorated table: input=%s", statement)
	}
	return decorated, nil
}

func addHour(targetTime time.Time, hadd float64) time.Time {
	if hadd == 0 {
		return targetTime
	}
	return targetTime.Add(time.Duration(hadd*60) * time.Minute)
}

func withHour(statement string, hour float64) (decorated string, err error) {
	beforeMSec := int(hour * 60 * 60 * 1000)

	// input:  tableName@
	// output: tableName@-3600000-
	decorated = tableRule.ReplaceAllString(statement, fmt.Sprintf("@-%d-", beforeMSec))

	if statement == decorated {
		return "", fmt.Errorf("Failed to decorated table: input=%s", statement)
	}
	return decorated, nil
}
