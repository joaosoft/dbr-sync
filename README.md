# monitor
[![Build Status](https://travis-ci.org/joaosoft/monitor.svg?branch=master)](https://travis-ci.org/joaosoft/monitor) | [![codecov](https://codecov.io/gh/joaosoft/monitor/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/monitor) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/monitor)](https://goreportcard.com/report/github.com/joaosoft/monitor) | [![GoDoc](https://godoc.org/github.com/joaosoft/monitor?status.svg)](https://godoc.org/github.com/joaosoft/monitor)

A simple and centralized process monitor

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Get process(es)
* Create process 
* Update process 
* Update process status
, Delete process(es)

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/monitor
```

## Usage 
This examples are available in the project at [monitor/examples](https://github.com/joaosoft/monitor/tree/master/examples)
```
import "github.com/joaosoft/monitor"

func main() {
	m, err := monitor.NewAuthentication()
	if err != nil {
		panic(err)
	}

	if err := m.Start(); err != nil {
		panic(err)
	}
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
