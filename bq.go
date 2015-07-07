package main

import (
	"os/exec"
)

func Query(statement string) string {
	args := buildArgs(statement)
	cmd := exec.Command("bq", args...)
	output, _ := cmd.CombinedOutput()
	return string(output)
}

func buildArgs(statement string) (args []string) {
	args = append(args, "query")

	if isDryRun {
		cflags := "--dry_run"
		args = append(args, cflags)
	}

	args = append(args, statement)
	return args
}
