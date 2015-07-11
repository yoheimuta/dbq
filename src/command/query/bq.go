package query

import (
	"os/exec"
)

// Bq is an Executor that run the bq command.
type Bq struct{}

// CreateBq initializes the Bq struct.
func CreateBq() *Bq {
	bq := &Bq{}
	return bq
}

// Query runs the bq query command.
func (b *Bq) Query(statement string) string {
	args := b.buildArgs(statement)
	cmd := exec.Command("bq", args...)
	output, _ := cmd.CombinedOutput()
	return string(output)
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
