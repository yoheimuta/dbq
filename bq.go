package main

import (
	"os/exec"
)

func Query(gflags string, cflags string, statement string) string {
	var args []string
	if gflags != "" {
		args = append(args, gflags)
	}

	args = append(args, "query")

	if cflags != "" {
		args = append(args, cflags)
	}

	args = append(args, statement)

	cmd := exec.Command("bq", args...)
	output, _ := cmd.CombinedOutput()
	return string(output)
}
