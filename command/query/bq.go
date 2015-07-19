package query

import (
	"fmt"
	"math"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/dustin/go-humanize"
)

var bytesRule = regexp.MustCompile("process (\\d*?) bytes of data")

// Bq is an Executor that run the bq command.
type Bq struct{}

// NewBq initializes the Bq struct.
func NewBq() *Bq {
	bq := &Bq{}
	return bq
}

var execCommand = exec.Command

// Query runs the bq query command.
func (b *Bq) Query(statement string) (out string) {
	args := b.buildArgs(statement)
	cmd := execCommand("bq", args...)
	output, _ := cmd.CombinedOutput()
	return string(output)
}

// GetHumanReadbleInfo changes the raw result to human readable one.
func (b *Bq) GetHumanReadbleInfo(result string) (info string) {
	matches := bytesRule.FindStringSubmatch(result)

	if len(matches) < 2 {
		return "Failed to change the result of dryRun to human readable one: result=" + result
	}

	bytes, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return "Failed to change the result of dryRun to human readable one: result=" + result + ": match=" + matches[1]
	}

	info = fmt.Sprintf("- %v bytes equal to %v bytes\n", bytes, humanize.Comma(bytes))
	info += fmt.Sprintf("- %v bytes equal to %v\n", bytes, humanize.IBytes(uint64(bytes)))

	tibibytes := float64(bytes) / math.Pow(1024, 4)
	costs := 5.0
	info += fmt.Sprintf("- %v bytes equal to $%.5f (= %.5f TiB * $%v)\n", bytes, tibibytes*costs, tibibytes, costs)
	return info
}

func (b *Bq) buildArgs(statement string) (args []string) {
	args = append(args, "query")

	if isDryRun {
		cflags := "--dry_run"
		args = append(args, cflags)
	}

	args = append(args, statement)
	return args
}
