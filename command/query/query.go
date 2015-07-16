package query

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var isVerbose bool
var isDryRun bool
var onlyStatement bool

// Args represents Arguments of the query command
type Args struct {
	beforeHour float64
	startDate  string
	endDate    string
	tz         float64
	buffer     float64
}

// Run is a facade method of the query command
func Run(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Println("Not Found a statement")
		return
	}
	statement := c.Args()[0]

	isVerbose = c.Bool("verbose")
	isDryRun = c.Bool("dryRun")
	onlyStatement = c.Bool("onlyStatement")

	args := Args{
		beforeHour: c.Float64("beforeHour"),
		startDate:  c.String("startDate"),
		endDate:    c.String("endDate"),
		tz:         c.Float64("tz"),
		buffer:     c.Float64("buffer"),
	}

	q := newQuery(statement, args)
	output, err := q.query()
	if err != nil {
		fmt.Printf("Failed to run the command\n:error=%v\n", err)
		return
	}
	if output != "" {
		fmt.Printf(output)
	}
}

// Query is an Implementation that decorates the statement and run the bq query
type Query struct {
	deco *Decorator
	bq   *Bq
}

var newQuery = func(statement string, args Args) *Query {
	return &Query{
		deco: NewDecorator(statement, args),
		bq:   NewBq(),
	}
}

func (q Query) query() (output string, err error) {
	if onlyStatement {
		return q.printStmt()
	}

	if isDryRun {
		return q.dryRun()
	}

	return q.run()
}

func (q Query) printStmt() (output string, err error) {
	dStmt, err := q.deco.Apply()
	if err != nil {
		return "", err
	}

	return dStmt, nil
}

func (q Query) dryRun() (output string, err error) {
	raw := q.deco.Revert()
	fmt.Printf("Raw: %v\n%v\n", raw, q.bq.Query(raw))

	dStmt, err := q.deco.Apply()
	if err != nil {
		return "", err
	}

	fmt.Printf("Decorated: %v\n", dStmt)
	return q.bq.Query(dStmt), nil
}

func (q Query) run() (output string, err error) {
	dStmt, err := q.deco.Apply()
	if err != nil {
		return "", err
	}

	fmt.Printf("Decorated: %v\n", dStmt)
	return q.bq.Query(dStmt), nil
}
