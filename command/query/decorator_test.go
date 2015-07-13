package query

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDecorator(t *testing.T) {
	Convey("CreateDecorator", t, func() {
		So(func() { CreateDecorator("", Args{}) }, ShouldNotPanic)
	})

	Convey("Apply", t, func() {

		Convey("When the statement doesn't include the placeholder of decorator", func() {
			statement := "SELECT * FROM [account.table]"

			d := CreateDecorator(statement, Args{startDate: "2015-07-08 17:00:00"})
			actual, err := d.Apply()

			So(actual, ShouldBeBlank)
			So(err, ShouldNotBeNil)
		})

		Convey("When the statement includes the placeholder of decorator", func() {
			statement := "SELECT * FROM [account.table@]"

			Convey("The argument includes the startDate", func() {

				Convey("The argument of startDate is formatted date string", func() {
					startDate := "2015-07-08 17:00:00"

					Convey("The argument of tz is 0, which means UTC", func() {
						d := CreateDecorator(statement, Args{startDate: startDate, buffer: 1})
						actual, err := d.Apply()

						expected := "SELECT * FROM [account.table@1436371200000-]"
						So(actual, ShouldEqual, expected)
						So(err, ShouldBeNil)
					})

					Convey("The argument of tz is -9, which means JST", func() {
						d := CreateDecorator(statement, Args{startDate: startDate, tz: -9, buffer: 1})
						actual, err := d.Apply()

						expected := "SELECT * FROM [account.table@1436338800000-]"
						So(actual, ShouldEqual, expected)
						So(err, ShouldBeNil)
					})
				})

				Convey("The argument of startDate is not formatted date string", func() {
					startDate := "2015/07/08 170000"

					d := CreateDecorator(statement, Args{startDate: startDate})
					actual, err := d.Apply()

					So(actual, ShouldBeBlank)
					So(err, ShouldNotBeNil)
				})
			})

			Convey("The argument includes the startDate and endDate", func() {

				Convey("The argument of startDate and endDate are formatted date string", func() {
					startDate := "2015-07-08 17:00:00"
					endDate := "2015-07-08 18:00:00"

					Convey("The argument of tz is 0, which means UTC", func() {
						d := CreateDecorator(statement, Args{startDate: startDate, endDate: endDate, buffer: 1})
						actual, err := d.Apply()

						expected := "SELECT * FROM [account.table@1436371200000-1436382000000]"
						So(actual, ShouldEqual, expected)
						So(err, ShouldBeNil)
					})

					Convey("The argument of tz is -9, which means JST", func() {
						d := CreateDecorator(statement, Args{startDate: startDate, endDate: endDate, tz: -9, buffer: 1})
						actual, err := d.Apply()

						expected := "SELECT * FROM [account.table@1436338800000-1436349600000]"
						So(actual, ShouldEqual, expected)
						So(err, ShouldBeNil)
					})
				})

				Convey("The argument of endDate is not formatted date string", func() {
					startDate := "2015-07-08 17:00:00"
					endDate := "2015/07/08 180000"

					d := CreateDecorator(statement, Args{startDate: startDate, endDate: endDate})
					actual, err := d.Apply()

					So(actual, ShouldBeBlank)
					So(err, ShouldNotBeNil)
				})
			})

			Convey("The argument includes the beforeHour", func() {

				Convey("The argument of beforeHour is greater than 0", func() {
					d := CreateDecorator(statement, Args{beforeHour: 3.0, buffer: 1})
					actual, err := d.Apply()

					expected := "SELECT * FROM [account.table@-10800000-]"
					So(actual, ShouldEqual, expected)
					So(err, ShouldBeNil)
				})

				Convey("The argument of beforeHour is 0", func() {
					d := CreateDecorator(statement, Args{beforeHour: 1.0, buffer: 1})
					actual, err := d.Apply()

					expected := "SELECT * FROM [account.table@-3600000-]"
					So(actual, ShouldEqual, expected)
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("When the statement includes the placeholders of decorator and customFunc", func() {
			statement := "SELECT * FROM [account.table@] WHERE _tz(2015-07-08 17:00:00) <= time"
			d := CreateDecorator(statement, Args{beforeHour: 3.0, tz: -9.0, buffer: 1})
			actual, err := d.Apply()

			expected := "SELECT * FROM [account.table@-10800000-] WHERE DATE_ADD('2015-07-08 17:00:00', -9, 'HOUR') <= time"
			So(actual, ShouldEqual, expected)
			So(err, ShouldBeNil)
		})
	})

	Convey("Revert", t, func() {
		statement := "SELECT * FROM [account.table@] WHERE _tz(2015-07-08 17:00:00) <= time ORDER BY time DESC"
		d := CreateDecorator(statement, Args{tz: -9})

		expected := "SELECT * FROM [account.table] WHERE DATE_ADD('2015-07-08 17:00:00', -9, 'HOUR') <= time ORDER BY time DESC"
		So(d.Revert(), ShouldEqual, expected)
	})
}
