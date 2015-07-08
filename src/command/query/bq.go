package query

import (
	"os/exec"
)

type Bq struct{}

func CreateBq() *Bq {
	bq := &Bq{}
	return bq
}

func (this *Bq) Query(statement string) string {
	args := this.buildArgs(statement)
	cmd := exec.Command("bq", args...)
	output, _ := cmd.CombinedOutput()
	return string(output)
}

func (this *Bq) buildArgs(statement string) (args []string) {
	args = append(args, "query")

	if isDryRun {
		cflags := "--dry_run"
		args = append(args, cflags)
	}

	args = append(args, statement)
	return args
}
