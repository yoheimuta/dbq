package query

import (
	"fmt"

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

func Run(c *cli.Context) {
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

	q := createQuery(statement, config)
	output, err := q.query()
	if err != nil {
		fmt.Printf("Failed to run the command\n: error=%v\n", err)
		return
	}
	if output != "" {
		fmt.Printf(output)
	}
}

type Query struct {
	deco *Decorator
	bq   *Bq
}

func createQuery(statement string, config Config) *Query {
	return &Query{
		deco: CreateDecorator(statement, config),
		bq:   CreateBq(),
	}
}

func (this Query) query() (output string, err error) {
	if onlyStatement {
		return this.printStmt()
	}

	if isDryRun {
		return this.dryRun()
	}

	return this.run()
}

func (this Query) printStmt() (output string, err error) {
	dStmt, err := this.deco.Apply()
	if err != nil {
		return "", err
	}

	return dStmt, nil
}

func (this Query) dryRun() (output string, err error) {
	raw := this.deco.Revert()
	fmt.Printf("Raw: %v\n%v\n", raw, this.bq.Query(raw))

	dStmt, err := this.deco.Apply()
	if err != nil {
		return "", err
	}

	fmt.Printf("Decorated: %v\n", dStmt)
	return this.bq.Query(dStmt), nil
}

func (this Query) run() (output string, err error) {
	dStmt, err := this.deco.Apply()
	if err != nil {
		return "", err
	}

	fmt.Printf("Decorated: %v\n", dStmt)
	return this.bq.Query(dStmt), nil
}
