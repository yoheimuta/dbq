package query

import (
	"fmt"
	"regexp"
	"strings"
)

var dateAddFuncRule = regexp.MustCompile("_tz\\((.*?)\\)")

type CustomFunc struct {
	config Config
}

func CreateCustomFunc(config Config) *CustomFunc {
	customFunc := &CustomFunc{
		config: config,
	}
	return customFunc
}

func (this *CustomFunc) Apply(statement string) string {
	return this._tz(statement)
}

func (this *CustomFunc) _tz(statement string) string {
	if this.config.hadd == 0 {
		return statement
	}

	replaced := statement
	matches := dateAddFuncRule.FindAllStringSubmatch(statement, -1)

	for _, match := range matches {
		if len(match) < 2 {
			if isVerbose {
				fmt.Printf("Skip to replace _tz(): statement=%v\n", statement)
			}
			return statement
		}

		date := match[1]
		old := fmt.Sprintf("_tz(%v)", date)
		new := fmt.Sprintf("DATE_ADD('%v', %v, 'HOUR')", date, this.config.hadd)
		replaced = strings.Replace(replaced, old, new, 1)
	}
	return replaced
}
