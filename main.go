package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var isVerbose bool
var isDryRun bool
var onlyStatement bool

func main() {
	app := cli.NewApp()
	app.Name = "dbq"
	app.Usage = "CLI tool to decorate bigquery table"
	app.Version = "0.0.1"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:      "query",
			ShortName: "q",
			Usage:     "Run bq query with complementing table decorator",
			Flags: []cli.Flag{
				cli.Float64Flag{
					Name:  "hour",
					Value: 0,
					Usage: "a decimal to specify the hour ago, relative to the current time",
				},
				cli.StringFlag{
					Name:  "start",
					Value: "",
					Usage: "a datetime to specify date range with end flag",
				},
				cli.StringFlag{
					Name:  "end",
					Value: "",
					Usage: "a datetime to specify date range with start flag",
				},
				cli.StringFlag{
					Name:  "hadd",
					Value: "",
					Usage: "+hour or -hour to add to start and end datetime, considering timezone",
				},
				cli.StringFlag{
					Name:  "gflags",
					Value: "",
					Usage: "[TODO] a gflags. see `bq help query`",
				},
				cli.StringFlag{
					Name:  "cflags",
					Value: "",
					Usage: "[TODO] a command flags. see `bq help query`",
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
			Action: func(c *cli.Context) {
				if len(c.Args()) != 1 {
					fmt.Println("Not Found a statement")
					return
				}

				statement := c.Args()[0]
				isVerbose = c.Bool("verbose")
				isDryRun = c.Bool("dryRun")
				onlyStatement = c.Bool("onlyStatement")

				output, err := query(statement, c.Float64("hour"), c.String("start"), c.String("end"), c.String("hadd"))
				if err != nil {
					fmt.Printf("Failed to run the command\n: error=%v\n", err)
					return
				}
				if output != "" {
					fmt.Printf(output)
				}
			},
		},
	}
	app.Run(os.Args)
}

func query(statement string, hour float64, start string, end string, hadd string) (output string, err error) {
	if isDryRun {
		dStatement := GetRawStatement(statement)
		fmt.Printf("Raw: %v\n", dStatement)
		dOutput := Query(dStatement)
		fmt.Printf("%v\n", dOutput)
	}

	decorated, err := Decorate(statement, hour, start, end, hadd)
	if err != nil {
		return "", err
	}
	if onlyStatement {
		fmt.Print(decorated)
		return "", nil
	}
	fmt.Printf("Decorated: %v\n", decorated)

	output = Query(decorated)
	return output, nil
}
