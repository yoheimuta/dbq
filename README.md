# dbq #

[![GitHub release](http://img.shields.io/github/release/yoheimuta/go-from-gist-to-issue.svg?style=flat-square)][release]
[![Wercker](http://img.shields.io/wercker/ci/54393fe184570fc622001411.svg?style=flat-square)][wercker]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[release]: https://github.com/yoheimuta/dbq/releases
[wercker]: https://app.wercker.com/project/bykey/d232745d4f00166be66404af31469036
[license]: https://github.com/yoheimuta/dbq/blob/master/LICENSE

CLI tool to easily Decorate BigQuery table name

## Description

`dbq` enables you to use [Table Range Decorators](https://cloud.google.com/bigquery/table-decorators) to perform a more cost-effective query to [BigQuery](https://cloud.google.com/bigquery/what-is-bigquery) without complex calculation.

- `dbq` supports both Relative value and Absolute value
- `dbq` also supports timezone calculation.

`dbq` will cut down on considerable data processed and spending :bowtie: :moneybag:

- [BigQuery offers on-demand pricing for queries](https://cloud.google.com/bigquery/pricing) and 1 TB of data processed costs $5.
- For example, a query to a 6TB table costs $30. Any SQL statements (ex. with LIMIT 1) don't affect the pricing.

## Installation

To install dbq, please use go get.

```
$ go get github.com/yoheimuta/dbq
...
$ dbq help
...
```

Or you can download a binary from [github relases page](https://github.com/yoheimuta/dbq/releases) and place it in $PATH directory.

## Requirements

- [bq-command-line-tool](https://cloud.google.com/bigquery/bq-command-line-tool)

## Usage

```ruby
# no-option equal to `--beforeHour=3`
$ dbq query "SELECT * FROM [foo.bar@]"

# equal to SELECT * FROM [foo.bar@-10800000-]
$ dbq query "SELECT * FROM [foo.bar@]" --beforeHour=3

# equal to SELECT * FROM [foo.bar@1436371200000-]
$ dbq query "SELECT * FROM [foo.bar@]" --startDate="2015-07-08 17:00:00"

# equal to SELECT * FROM [foo.bar@1436371200000-1436382000000]
$ dbq query "SELECT * FROM [foo.bar@]" --startDate="2015-07-08 17:00:00" --endDate="2015-07-08 18:00:00"

# equal to SELECT * FROM [foo.bar@1436338800000-]
$ dbq query "SELECT * FROM [foo.bar@]" --startDate="2015-07-08 17:00:00" --tz="-9"

# equal to SELECT * FROM [foo.bar@1436338800000-] WHERE DATE_ADD('2015-07-08 17:00:00', -9, 'HOUR') <= time and time <= DATE_ADD('2015-07-08 18:00:00', -9, 'HOUR')
$ dbq query "SELECT * FROM [foo.bar@] WHERE _tz(2015-07-08 17:00:00) <= time and time <= _tz(2015-07-08 18:00:00)" --startDate="2015-07-08 17:00:00" --tz="-9"
```

### Placeholders

- `@` will be replaced with `@<time1>-<time2>`
 - required
- `_tz(datetime)` will be replaced with `DATE_ADD('datetime', tz value, 'HOUR')`
 - optional

## DryRun

The option of `dryRun` shows how much cut down full scan bytes, so I strongly recommend to use this option before running any queries.

- A query with no table decorator will process 6.5 TB, then costs `6.5 * $5 = $32.5`.
- A query with table decorator will process 4.2 GB, then costs `0.0042 * $5 = $0.021`. `dbq` will save `$32.479`.

```ruby
$ dbq query "SELECT * FROM [foo.bar@]" --dryRun
Raw: SELECT * FROM [foo.bar]
Query successfully validated. Assuming the tables are not modified, running this query will process 6488095769102 bytes of data.

Decorated: SELECT * FROM [foo.bar@-10800000-]
Query successfully validated. Assuming the tables are not modified, running this query will process 4201020576 bytes of data.
```

## Options

```
$ dbq help query
NAME:
   query - Run bq query with complementing table range decorator

USAGE:
   command query [command options] [arguments...]

DESCRIPTION:


OPTIONS:
   --beforeHour '3'     a decimal to specify the hour ago, relative to the current time
   --startDate          a datetime to specify date range with end flag
   --endDate            a datetime to specify date range with start flag
   --tz '0'             a decimal of hour or -hour to add to start and end datetime, considering timezone
   --buffer '1'         a decimal of hour to add to start and end datetime, it's heuristic value
   --gflags             no support. Use onlyStatement instead
   --cflags             no support. Use onlyStatement instead
   --verbose            a flag to output verbosely
   --dryRun             a flag to run without any changes
   --onlyStatement      a flag to output only a decorated statement
```

## CHANGELOG

See [CHANGELOG](CHANGELOG.md)
