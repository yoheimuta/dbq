package main

import (
	"fmt"
	"regexp"
	"time"
)

var tableRule = regexp.MustCompile("@")

type Decorator struct {
	statement  string
	config     Config
	customFunc *CustomFunc
}

func CreateDecorator(statement string, config Config) *Decorator {
	cf := CreateCustomFunc(config)

	decorator := &Decorator{
		statement:  statement,
		config:     config,
		customFunc: cf,
	}
	return decorator
}

func (this *Decorator) Apply() (decorated string, err error) {
	stmt := this.customFunc.Apply(this.statement)

	if this.config.startDate != "" {
		startMSec, endMSec, err := this.getRangeMSec()
		if err != nil {
			return "", err
		}
		decorated = this.useAbs(stmt, startMSec, endMSec)
	} else {
		hour := this.config.hour
		if hour <= 0 {
			hour = 1.0
		}
		beforeMSec := int(hour * 60 * 60 * 1000)
		decorated = this.useRel(stmt, beforeMSec)
	}

	if stmt == decorated {
		return "", fmt.Errorf("Failed to decorated table: input=%s", stmt)
	}
	return decorated, nil
}

func (this *Decorator) Revert() (raw string) {
	stmt := this.customFunc.Apply(this.statement)

	// input:  tableName@
	// output: tableName
	raw = tableRule.ReplaceAllString(stmt, "")
	return raw
}

func (this *Decorator) getRangeMSec() (startMSec int64, endMSec int64, err error) {
	hadd := this.config.hadd
	buffer := this.config.buffer

	startTime, err := time.Parse("2006-01-02 15:04:05", this.config.startDate)
	if err != nil {
		return 0, 0, err
	}
	startTime = startTime.Add(time.Duration((hadd+buffer*-1)*60) * time.Minute)
	startMSec = startTime.Unix() * 1000

	end := this.config.endDate
	if end == "" {
		return startMSec, 0, nil
	} else {
		endTime, err := time.Parse("2006-01-02 15:04:05", end)
		if err != nil {
			return 0, 0, err
		}
		endTime = endTime.Add(time.Duration((hadd+buffer)*60) * time.Minute)
		endMSec = endTime.Unix() * 1000
		return startMSec, endMSec, nil
	}
}

func (this *Decorator) useAbs(statement string, startMSec int64, endMSec int64) string {
	var replaced string
	if endMSec == 0 {
		replaced = fmt.Sprintf("@%d-", startMSec)
	} else {
		replaced = fmt.Sprintf("@%d-%d", startMSec, endMSec)
	}

	// input:  tableName@
	// output: tableName@1435997334864-1436001613000
	decorated := tableRule.ReplaceAllString(statement, replaced)
	return decorated
}

func (this *Decorator) useRel(statement string, beforeMSec int) string {
	// input:  tableName@
	// output: tableName@-3600000-
	decorated := tableRule.ReplaceAllString(statement, fmt.Sprintf("@-%d-", beforeMSec))
	return decorated
}
