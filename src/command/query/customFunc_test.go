package query

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomFunc(t *testing.T) {
	Convey("CreateCustomFunc", t, func() {
		So(func() { CreateCustomFunc(Config{}) }, ShouldNotPanic)
	})
}
