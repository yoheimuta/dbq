package query

import (
	"fmt"
	"regexp"
	"strings"
)

var dateAddFuncRule = regexp.MustCompile("_tz\\((.*?)\\)")

// CustomFunc handles Custom Functions that are applied to the statement
type CustomFunc struct {
	args Args
}

// NewCustomFunc initializes the CustomFunc struct.
func NewCustomFunc(args Args) *CustomFunc {
	customFunc := &CustomFunc{
		args: args,
	}
	return customFunc
}

// Apply is a facade method that runs custom functions
func (c *CustomFunc) Apply(statement string) string {
	return c._tz(statement)
}

func (c *CustomFunc) _tz(statement string) string {
	if c.args.tz == 0 {
		if isVerbose {
			fmt.Printf("Skip to replace _tz(): tz=0\n")
		}
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
		new := fmt.Sprintf("DATE_ADD('%v', %v, 'HOUR')", date, c.args.tz)
		replaced = strings.Replace(replaced, old, new, 1)
	}
	return replaced
}
