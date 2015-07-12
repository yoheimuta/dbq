# dbq #

[![GitHub release](http://img.shields.io/github/release/yoheimuta/go-from-gist-to-issue.svg?style=flat-square)][release]
[![Wercker](http://img.shields.io/wercker/ci/54393fe184570fc622001411.svg?style=flat-square)][wercker]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[release]: https://github.com/yoheimuta/dbq/releases
[wercker]: https://app.wercker.com/project/bykey/d232745d4f00166be66404af3146903
[license]: https://github.com/yoheimuta/dbq/blob/master/LICENSE

CLI tool to easily decorate bigquery table name

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

```
$ dbq query "SELECT col1, col2 FROM [foo.bar@]"
$ dbq query "SELECT col1, col2 FROM [foo.bar@]" --hour=3
$ dbq query "SELECT col1, col2 FROM [foo.bar@]" --start="2015-07-08 17:00:00"
$ dbq query "SELECT col1, col2 FROM [foo.bar@]" --start="2015-07-08 17:00:00" --end="2015-07-08 18:00:00"
$ dbq query "SELECT col1, col2 FROM [foo.bar@]" --hadd="-9"
$ dbq query "SELECT col1, col2 FROM [foo.bar@] WHERE _tz(2015-07-08 17:00:00) <= time and time <= _tz(2015-07-08 18:00:00)" --hadd="-9"
```

## Options

```
$ dbq help query
NAME:
   query - Run bq query with complementing table decorator

USAGE:
   command query [command options] [arguments...]

DESCRIPTION:


OPTIONS:
   --hour '0'           a decimal to specify the hour ago, relative to the current time
   --start              a datetime to specify date range with end flag
   --end                a datetime to specify date range with start flag
   --hadd '0'           a decimal of hour or -hour to add to start and end datetime, considering timezone
   --buffer '1'         a decimal of hour to add to start and end datetime, it's heuristic value
   --gflags             no support. Use onlyStatement instead
   --cflags             no support. Use onlyStatement instead
   --verbose            a flag to output verbosely
   --dryRun             a flag to run without any changes
   --onlyStatement      a flag to output only a decorated statement
```

## CHANGELOG

See [CHANGELOG](CHANGELOG.md)
