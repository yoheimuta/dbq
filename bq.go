package main

import (
	"os/exec"
	"regexp"
)

var dryRunRule = regexp.MustCompile("--dry_run")

func Query(gflags string, cflags string, statement string) string {
	args := buildArgs(gflags, cflags, statement)
	cmd := exec.Command("bq", args...)
	output, _ := cmd.CombinedOutput()
	return string(output)
}

func buildArgs(gflags, cflags string, statement string) (args []string) {
	if gflags != "" {
		args = append(args, gflags)
	}

	args = append(args, "query")

	if isDryRun {
		if cflags == "" {
			cflags = "--dry_run"
		} else if !dryRunRule.MatchString(cflags) {
			cflags += " --dry_run"
		}
	}
	if cflags != "" {
		args = append(args, cflags)
	}

	args = append(args, statement)
	return args
}
