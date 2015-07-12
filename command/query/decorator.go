package query

import (
	"fmt"
	"regexp"
	"time"
)

var tableRule = regexp.MustCompile("@")

// Decorator transforms the no-decorated statement into the decorated one
type Decorator struct {
	statement  string
	args       Args
	customFunc *CustomFunc
}

// CreateDecorator initializes the Decorator struct.
func CreateDecorator(statement string, args Args) *Decorator {
	cf := CreateCustomFunc(args)

	decorator := &Decorator{
		statement:  statement,
		args:       args,
		customFunc: cf,
	}
	return decorator
}

// Apply is a facade method that transforms the no-decorated statement into the decorated one
func (d *Decorator) Apply() (decorated string, err error) {
	stmt := d.customFunc.Apply(d.statement)

	if d.args.startDate != "" {
		startMSec, endMSec, err := d.getRangeMSec()
		if err != nil {
			return "", err
		}
		decorated = d.useAbs(stmt, startMSec, endMSec)
	} else {
		hour := d.args.hour
		if hour <= 0 {
			hour = 1.0
		}
		beforeMSec := int(hour * 60 * 60 * 1000)
		decorated = d.useRel(stmt, beforeMSec)
	}

	if stmt == decorated {
		return "", fmt.Errorf("Failed to decorated table: input=%s", stmt)
	}
	return decorated, nil
}

// Revert transforms the statement with @ into the one without @
func (d *Decorator) Revert() (raw string) {
	stmt := d.customFunc.Apply(d.statement)

	// input:  tableName@
	// output: tableName
	raw = tableRule.ReplaceAllString(stmt, "")
	return raw
}

func (d *Decorator) getRangeMSec() (startMSec int64, endMSec int64, err error) {
	hadd := d.args.hadd
	buffer := d.args.buffer

	// startTime: convert the formatted string into the unixtime ms
	startTime, err := time.Parse("2006-01-02 15:04:05", d.args.startDate)
	if err != nil {
		return 0, 0, err
	}
	startTime = startTime.Add(time.Duration((hadd+buffer*-1)*60) * time.Minute)
	startMSec = startTime.Unix() * 1000

	// endTime: convert the formatted string into the unixtime ms
	end := d.args.endDate
	if end == "" {
		return startMSec, 0, nil
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", end)
	if err != nil {
		return 0, 0, err
	}
	endTime = endTime.Add(time.Duration((hadd+buffer)*60) * time.Minute)
	endMSec = endTime.Unix() * 1000
	return startMSec, endMSec, nil
}

func (d *Decorator) useAbs(statement string, startMSec int64, endMSec int64) string {
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

func (d *Decorator) useRel(statement string, beforeMSec int) string {
	// input:  tableName@
	// output: tableName@-3600000-
	decorated := tableRule.ReplaceAllString(statement, fmt.Sprintf("@-%d-", beforeMSec))
	return decorated
}
