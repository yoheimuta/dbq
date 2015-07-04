package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var tableRule = regexp.MustCompile("@")
var haddRule = regexp.MustCompile("([+|-])(\\d)")

func Decorate(query string, hour float64, start string, end string, hadd string) (decorated string, err error) {
	if start != "" {
		decorated, err = withDateTime(query, start, end, hadd)
	} else {
		if hour <= 0 {
			hour = 0.5
		}
		decorated, err = withHour(query, hour)
	}
	return decorated, err
}

func withDateTime(query string, start string, end string, hadd string) (decorated string, err error) {
	startTime, err := time.Parse("2006-01-02 15:04:05", start)
	if err != nil {
		return "", err
	}
	startTime = hourAdd(startTime, hadd)
	startMSec := startTime.Unix() * 1000

	var replaced string
	if end == "" {
		replaced = fmt.Sprintf("@%d-", startMSec)
	} else {
		endTime, err := time.Parse("2006-01-02 15:04:05", end)
		if err != nil {
			return "", err
		}
		endTime = hourAdd(endTime, hadd)
		endMSec := endTime.Unix() * 1000
		replaced = fmt.Sprintf("@%d-%d", startMSec, endMSec)
	}

	// input:  tableName@
	// output: tableName@1435997334864-1436001613000
	decorated = tableRule.ReplaceAllString(query, replaced)

	if query == decorated {
		return "", fmt.Errorf("Failed to decorated table: input=%s", query)
	}
	return decorated, nil
}

func hourAdd(targetTime time.Time, hadd string) time.Time {
	if hadd == "" {
		return targetTime
	}
	matches := haddRule.FindAllStringSubmatch(hadd, -1)

	if len(matches) < 1 {
		if isVerbose {
			fmt.Printf("Skip to add hadd hour because no match: hadd=%v\n", hadd)
		}
		return targetTime
	}
	if len(matches[0]) < 3 {
		if isVerbose {
			fmt.Printf("Skip to add hadd hour because invalid format: hadd=%v\n", hadd)
		}
		return targetTime
	}

	diff, _ := strconv.Atoi(matches[0][2])
	if matches[0][1] == "-" {
		diff *= -1
	}
	return targetTime.Add(time.Duration(diff) * time.Hour)
}

func withHour(query string, hour float64) (decorated string, err error) {
	beforeMSec := int(hour * 60 * 60 * 1000)

	// input:  tableName@
	// output: tableName@-3600000-
	decorated = tableRule.ReplaceAllString(query, fmt.Sprintf("@-%d-", beforeMSec))

	if query == decorated {
		return "", fmt.Errorf("Failed to decorated table: input=%s", query)
	}
	return decorated, nil
}
