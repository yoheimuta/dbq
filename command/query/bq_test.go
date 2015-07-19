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

	Convey("When the result of dryRun query is changed to humanreadable one", t, func() {
		b := NewBq()
		result := "Query successfully validated. Assuming the tables are not modified, running this query will process 8133291239 bytes of data."

		info := b.GetHumanReadbleInfo(result)

		line1 := "- 8133291239 bytes equal to 8,133,291,239 bytes"
		line2 := "- 8133291239 bytes equal to 7.6GiB"
		line3 := "- 8133291239 bytes equal to $0.03699 (= 0.00740 TiB * $5)"
		So(info, ShouldEqual, fmt.Sprintf("%v\n%v\n%v\n", line1, line2, line3))
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
