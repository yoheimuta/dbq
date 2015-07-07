package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var isVerbose bool
var isDryRun bool
var onlyStatement bool

type Config struct {
	hour      float64
	startDate string
	endDate   string
	hadd      float64
	buffer    float64
}

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
				cli.Float64Flag{
					Name:  "hadd",
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
			Action: func(c *cli.Context) {
				if len(c.Args()) != 1 {
					fmt.Println("Not Found a statement")
					return
				}
				statement := c.Args()[0]

				isVerbose = c.Bool("verbose")
				isDryRun = c.Bool("dryRun")
				onlyStatement = c.Bool("onlyStatement")

				config := Config{
					hour:      c.Float64("hour"),
					startDate: c.String("start"),
					endDate:   c.String("end"),
					hadd:      c.Float64("hadd"),
					buffer:    c.Float64("buffer"),
				}

				output, err := query(statement, config)
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

func query(statement string, config Config) (output string, err error) {
	deco := CreateDecorator(statement, config)
	bq := CreateBq()

	if onlyStatement {
		return printStmt(deco)
	}

	if isDryRun {
		return dryRun(deco, bq)
	}

	return run(deco, bq)
}

func printStmt(deco *Decorator) (output string, err error) {
	dStmt, err := deco.Apply()
	if err != nil {
		return "", err
	}

	return dStmt, nil
}

func dryRun(deco *Decorator, bq *Bq) (output string, err error) {
	raw := deco.Revert()
	fmt.Printf("Raw: %v\n%v\n", raw, bq.Query(raw))

	dStmt, err := deco.Apply()
	if err != nil {
		return "", err
	}

	fmt.Printf("Decorated: %v\n", dStmt)
	return bq.Query(dStmt), nil
}

func run(deco *Decorator, bq *Bq) (output string, err error) {
	dStmt, err := deco.Apply()
	if err != nil {
		return "", err
	}

	fmt.Printf("Decorated: %v\n", dStmt)
	return bq.Query(dStmt), nil
}
