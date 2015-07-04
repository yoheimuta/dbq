package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var isVerbose bool

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
			Usage:     "same to bq query",
			Flags: []cli.Flag{
				cli.Float64Flag{
					Name:  "hour",
					Value: 0,
					Usage: "a decimal to decorate hour ago, relative to the current time",
				},
				cli.StringFlag{
					Name:  "start",
					Value: "",
					Usage: "a datetime to specify date range",
				},
				cli.StringFlag{
					Name:  "end",
					Value: "",
					Usage: "a datetime to specify date range",
				},
				cli.StringFlag{
					Name:  "hadd",
					Value: "",
					Usage: "an hour to add start and end datetime, because of timezone",
				},
				cli.BoolFlag{
					Name:  "verbose",
					Usage: "a flag to log verbosely",
				},
			},
			Action: func(c *cli.Context) {
				if len(c.Args()) != 1 {
					fmt.Println("Not Found a query")
					return
				}

				query := c.Args()[0]
				isVerbose = c.Bool("verbose")

				output, err := action(query, c.Float64("hour"), c.String("start"), c.String("end"), c.String("hadd"))
				if err != nil {
					fmt.Printf("Failed to run the command: %v\n", err)
					return
				}
				fmt.Printf(output)
			},
		},
	}
	app.Run(os.Args)
}

func action(query string, hour float64, start string, end string, hadd string) (output string, err error) {
	decorated, err := Decorate(query, hour, start, end, hadd)
	if err != nil {
		return "", err
	}
	fmt.Printf("decorated: %v\n", decorated)

	output = Query("", "", decorated)
	return output, nil
}
