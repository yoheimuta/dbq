package query

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomFunc(t *testing.T) {
	Convey("CreateCustomFunc", t, func() {
		So(func() { CreateCustomFunc(Args{}) }, ShouldNotPanic)
	})

	Convey("When the statement includes the placeholder of _tz", t, func() {

		Convey("When the placeholder of _tz appears once", func() {
			statement := "SELECT * FROM [account.table@] WHERE col='val' and _tz(2015-07-08 17:00:00) <= time ORDER BY time DESC"

			Convey("The argument of tz is 0, which means UTC", func() {
				cf := CreateCustomFunc(Args{tz: 0})
				So(cf.Apply(statement), ShouldEqual, statement)
			})

			Convey("The argument of tz is -9, which means JST", func() {
				cf := CreateCustomFunc(Args{tz: -9})
				So(cf.Apply(statement), ShouldNotEqual, statement)

				expected := "SELECT * FROM [account.table@] WHERE col='val' and DATE_ADD('2015-07-08 17:00:00', -9, 'HOUR') <= time ORDER BY time DESC"
				So(cf.Apply(statement), ShouldEqual, expected)
			})
		})

		Convey("When the placeholder of _tz appears twice", func() {
			statement := "SELECT * FROM [account.table@] WHERE col='val' and _tz(2015-07-08 17:00:00) <= time and time <= _tz(2015-07-08 18:00:00) ORDER BY time DESC"

			Convey("The argument of tz is 0, which means UTC", func() {
				cf := CreateCustomFunc(Args{tz: 0})
				So(cf.Apply(statement), ShouldEqual, statement)
			})

			Convey("The argument of tz is -9, which means JST", func() {
				cf := CreateCustomFunc(Args{tz: -9})
				So(cf.Apply(statement), ShouldNotEqual, statement)

				expected := "SELECT * FROM [account.table@] WHERE col='val' and DATE_ADD('2015-07-08 17:00:00', -9, 'HOUR') <= time and time <= DATE_ADD('2015-07-08 18:00:00', -9, 'HOUR') ORDER BY time DESC"
				So(cf.Apply(statement), ShouldEqual, expected)
			})
		})
	})

	Convey("When the statement doesn't include placeholder of _tz", t, func() {
		statement := "SELECT * FROM [account.table@] WHERE col='val' and DATE_ADD('2015-07-08 17:00:00', -9, 'HOUR') ORDER BY time DESC"
		cf := CreateCustomFunc(Args{tz: -9})
		So(cf.Apply(statement), ShouldEqual, statement)
	})
}
