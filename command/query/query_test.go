package query

import (
	"flag"
	"os/exec"
	"regexp"
	"testing"

	"github.com/codegangsta/cli"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRun(t *testing.T) {
	Convey("When the arguements are parsed", t, func() {
		expectedStmt := "SELECT * FROM [account.table]"
		expectedBeforeHour := 1.0
		expectedStart := "2015-07-08 17:00:00"
		expectedEnd := "2015-07-08 18:00:00"
		expectedHadd := -9.0
		expectedBuffer := 1.0

		prevNewQuery := newQuery
		newQuery = func(statement string, args Args) *Query {
			So(statement, ShouldEqual, expectedStmt)

			So(args.beforeHour, ShouldEqual, expectedBeforeHour)
			So(args.startDate, ShouldEqual, expectedStart)
			So(args.endDate, ShouldEqual, expectedEnd)
			So(args.tz, ShouldEqual, expectedHadd)
			So(args.buffer, ShouldEqual, expectedBuffer)
			So(isVerbose, ShouldBeFalse)
			So(isDryRun, ShouldBeTrue)
			So(onlyStatement, ShouldBeTrue)

			return &Query{
				deco: NewDecorator(statement, args),
				bq:   NewBq(),
			}
		}
		defer func() { newQuery = prevNewQuery }()

		set := flag.NewFlagSet("test", 0)
		set.Parse([]string{expectedStmt})
		set.Float64("beforeHour", expectedBeforeHour, "")
		set.String("startDate", expectedStart, "")
		set.String("endDate", expectedEnd, "")
		set.Float64("tz", expectedHadd, "")
		set.Float64("buffer", expectedBuffer, "")
		set.Bool("verbose", false, "")
		set.Bool("dryRun", true, "")
		set.Bool("onlyStatement", true, "")
		c := cli.NewContext(nil, set, nil)

		Run(c)
	})
}

func TestQuery(t *testing.T) {
	statement := "SELECT * FROM [account.table@] WHERE _tz(2015-07-08 17:00:00) <= time"
	args := Args{
		beforeHour: 3.0,
		tz:         -9.0,
		buffer:     1.0,
	}
	q := newQuery(statement, args)

	Convey("When the onlyStatement flag is on", t, func() {
		onlyStatement = true
		actual, err := q.query()

		expected := "SELECT * FROM [account.table@-10800000-] WHERE DATE_ADD('2015-07-08 17:00:00', -9, 'HOUR') <= time"
		So(actual, ShouldEqual, expected)
		So(err, ShouldBeNil)
	})

	Convey("When the isDryRun flag is on", t, func() {
		execCommand = fakeExecCommand
		defer func() { execCommand = exec.Command }()

		onlyStatement = false
		isDryRun = true
		actual, err := q.query()

		So(actual, ShouldEqual, "")
		So(err, ShouldBeNil)
	})

	Convey("When the onlyStatement and isDryRun flag are off", t, func() {
		execCommand = fakeExecCommand
		defer func() { execCommand = exec.Command }()

		onlyStatement = false
		isDryRun = false
		actual, err := q.query()

		So(regexp.MustCompile(mockReturnVal).MatchString(actual), ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}
