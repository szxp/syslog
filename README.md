[![Build Status](https://travis-ci.org/szxp/syslog.svg?branch=master)](https://travis-ci.org/szxp/syslog)
[![Build Status](https://ci.appveyor.com/api/projects/status/github/szxp/syslog?branch=master&svg=true)](https://ci.appveyor.com/project/szxp/syslog)
[![GoDoc](https://godoc.org/github.com/szxp/syslog?status.svg)](https://godoc.org/github.com/szxp/syslog)
[![Go Report Card](https://goreportcard.com/badge/github.com/szxp/syslog)](https://goreportcard.com/report/github.com/szxp/syslog)

# syslog
Syslog message formatter.


## Example
```go
package main

import (
	"github.com/szxp/syslog"
	"log"
	"os"
)

func main() {
	const msg = "Start HTTP server (addr=:8080)"

	wrappedWriter := syslog.NewWriter(os.Stdout, syslog.USER|syslog.NOTICE)
	logger := log.New(wrappedWriter, "", 0)
	logger.Println(msg)

	// Output is similar to this:
	// <13>1 2017-08-15T23:13:15.33+02:00 laptop /path/to/myprogram 21650 - - Start HTTP server (addr=:8080)
}
```


