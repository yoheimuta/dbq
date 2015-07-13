package main

import (
	"os"

	"github.com/yoheimuta/dbq/command/query"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "dbq"
	app.Usage = "CLI tool to decorate bigquery table"
	app.Version = "0.0.2"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:      "query",
			ShortName: "q",
			Usage:     "Run bq query with complementing table range decorator",
			Flags: []cli.Flag{
				cli.Float64Flag{
					Name:  "beforeHour",
					Value: 3,
					Usage: "a decimal to specify the hour ago, relative to the current time",
				},
				cli.StringFlag{
					Name:  "startDate",
					Value: "",
					Usage: "a datetime to specify date range with end flag",
				},
				cli.StringFlag{
					Name:  "endDate",
					Value: "",
					Usage: "a datetime to specify date range with start flag",
				},
				cli.Float64Flag{
					Name:  "tz",
					Value: 0,
					Usage: "a decimal of hour or -hour to add to start and end datetime, considering timezone",
				},
				cli.Float64Flag{
					Name:  "buffer",
					Value: 1,
					Usage: "a decimal of hour to add to start and end datetime, it's heuristic value",
				},
				cli.StringFlag{
					Name:  "gflags",
					Value: "",
					Usage: "no support. Use onlyStatement instead",
				},
				cli.StringFlag{
					Name:  "cflags",
					Value: "",
					Usage: "no support. Use onlyStatement instead",
				},
				cli.BoolFlag{
					Name:  "verbose",
					Usage: "a flag to output verbosely",
				},
				cli.BoolFlag{
					Name:  "dryRun",
					Usage: "a flag to run without any changes",
				},
				cli.BoolFlag{
					Name:  "onlyStatement",
					Usage: "a flag to output only a decorated statement",
				},
			},
			Action: query.Run,
		},
	}
	app.Run(os.Args)
}
