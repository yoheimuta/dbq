package query

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

const mockReturnVal = "RESULT"

func TestBq(t *testing.T) {
	Convey("NewBq", t, func() {
		So(func() { NewBq() }, ShouldNotPanic)
	})

	Convey("When the bq command executes the statement", t, func() {
		execCommand = fakeExecCommand
		defer func() { execCommand = exec.Command }()

		actual := NewBq().Query("SELECT * FROM [account.table]")
		So(regexp.MustCompile(mockReturnVal).MatchString(actual), ShouldBeTrue)
	})

	Convey("When the arguments are built", t, func() {
		b := NewBq()
		statement := "SELECT * FROM [account.table]"

		Convey("When the dryRun flag is off", func() {
			actual := b.buildArgs(statement)
			So(actual[0], ShouldEqual, "query")
			So(actual[1], ShouldEqual, statement)
		})

		Convey("When the dryRun flag is on", func() {
			isDryRun = true
			actual := b.buildArgs(statement)
			So(actual[0], ShouldEqual, "query")
			So(actual[1], ShouldEqual, "--dry_run")
			So(actual[2], ShouldEqual, statement)
			isDryRun = false
		})
	})
}

// TestHelperProcess isn't a real test. It's used as a helper process
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	fmt.Fprintf(os.Stdout, mockReturnVal)
}
