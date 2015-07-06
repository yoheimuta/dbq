package main

import (
	"fmt"
	"regexp"
	"strings"
)

var dateAddFuncRule = regexp.MustCompile("_ha\\((.*?)\\)")

func DateAdd(statement string, hadd float64) (replaced string, err error) {
	if hadd == 0 {
		return statement, nil
	}

	replaced = statement
	matches := dateAddFuncRule.FindAllStringSubmatch(statement, -1)

	for _, match := range matches {
		if len(match) < 2 {
			if isVerbose {
				fmt.Printf("Skip to replace ha because invalid format: statement=%v\n", statement)
			}
			return statement, nil
		}

		date := match[1]
		old := fmt.Sprintf("_ha(%v)", date)
		new := fmt.Sprintf("DATE_ADD('%v', %v, 'HOUR')", date, hadd)
		replaced = strings.Replace(replaced, old, new, 1)
	}
	return replaced, nil
}
