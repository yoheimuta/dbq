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
					Value: 0.5,
					Usage: "a decimal to decorate hour ago, relative to the current time",
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

				hour := c.Float64("hour")
				isVerbose = c.Bool("verbose")

				var output string
				var err error
				if 0 < hour {
					output, err = action(query, hour)
				}

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

func action(query string, hour float64) (output string, err error) {
	if isVerbose {
		fmt.Printf("query: %v, hour: %v\n", query, hour)
	}

	decorated, err := Decorate(hour, query)
	if err != nil {
		return "", err
	}
	fmt.Printf("decorated: %v\n", decorated)

	output = Query("", "", decorated)
	return output, nil
}
