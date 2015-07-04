package main

import (
	"fmt"
	"regexp"
)

var tableRule = regexp.MustCompile("@")

func Decorate(hour float64, query string) (decorated string, err error) {
	beforeMSec := int(hour * 60 * 60 * 1000)

	// input:  tableName@
	// output: tableName@-3600000-
	decorated = tableRule.ReplaceAllString(query, fmt.Sprintf("@-%d-", beforeMSec))

	if query == decorated {
		return "", fmt.Errorf("Failed to decorated table: input=%s", query)
	}
	return decorated, nil
}
