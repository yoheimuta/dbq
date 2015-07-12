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
		expectedHour := 1.0
		expectedStart := "2015-07-08 17:00:00"
		expectedEnd := "2015-07-08 18:00:00"
		expectedHadd := -9.0
		expectedBuffer := 1.0

		prevCreateQuery := createQuery
		createQuery = func(statement string, args Args) *Query {
			So(statement, ShouldEqual, expectedStmt)

			So(args.hour, ShouldEqual, expectedHour)
			So(args.startDate, ShouldEqual, expectedStart)
			So(args.endDate, ShouldEqual, expectedEnd)
			So(args.hadd, ShouldEqual, expectedHadd)
			So(args.buffer, ShouldEqual, expectedBuffer)
			So(isVerbose, ShouldBeFalse)
			So(isDryRun, ShouldBeTrue)
			So(onlyStatement, ShouldBeTrue)

			return &Query{
				deco: CreateDecorator(statement, args),
				bq:   CreateBq(),
			}
		}
		defer func() { createQuery = prevCreateQuery }()

		set := flag.NewFlagSet("test", 0)
		set.Parse([]string{expectedStmt})
		set.Float64("hour", expectedHour, "")
		set.String("start", expectedStart, "")
		set.String("end", expectedEnd, "")
		set.Float64("hadd", expectedHadd, "")
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
		hour:   3.0,
		hadd:   -9.0,
		buffer: 1.0,
	}
	q := createQuery(statement, args)

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

		So(regexp.MustCompile(mockReturnVal).MatchString(actual), ShouldBeTrue)
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
